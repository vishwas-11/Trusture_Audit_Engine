import { useEffect, useRef, useState } from "react";
import api from "../../api/axios";

type Donation = {
  id: string;
  tx_hash: string;
  donor_wallet: string;
  ngo_wallet: string;
  purpose: string;
  amount: number;
  status: string;
  created_at: number;
};

export default function NGODashboard() {
  const [walletInput, setWalletInput] = useState<string>("");
  const [walletAddress, setWalletAddress] = useState<string>("");
  const [donations, setDonations] = useState<Donation[]>([]);
  const [selectedDonation, setSelectedDonation] = useState<string>("");
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState<string>("");
  const [loadingTransactions, setLoadingTransactions] = useState(false);
  const [transactionsError, setTransactionsError] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  useEffect(() => {
    if (!walletAddress) return;
    setLoadingTransactions(true);
    setTransactionsError("");

    api
      .get("/ngo/dashboard", {
        headers: {
          "X-ROLE": "NGO",
          "X-WALLET": walletAddress,
        },
      })
      .then((res) => setDonations(res.data))
      .catch((err) => {
        console.error(err);
        setDonations([]);
        setTransactionsError("Failed to fetch transactions for this wallet.");
      })
      .finally(() => setLoadingTransactions(false));
  }, [walletAddress]);

  const loadTransactions = () => {
    const trimmedWallet = walletInput.trim();
    if (!/^0x[a-fA-F0-9]{40}$/.test(trimmedWallet)) {
      setTransactionsError("Enter a valid MetaMask wallet address.");
      return;
    }

    setWalletAddress(trimmedWallet);
    setSelectedDonation("");
    setStatus("");
  };

  const uploadProof = async () => {
    if (!walletAddress) {
      alert("Enter NGO wallet address first");
      return;
    }

    if (!file || !selectedDonation) {
      alert("Select donation and file");
      return;
    }

    const supportedMimeTypes = new Set([
      "application/pdf",
      "image/jpeg",
      "image/png",
      "image/webp",
    ]);
    const supportedExtensions = [".pdf", ".jpg", ".jpeg", ".png", ".webp"];
    const hasValidExtension = supportedExtensions.some((ext) =>
      file.name.toLowerCase().endsWith(ext),
    );
    if (!supportedMimeTypes.has(file.type) && !hasValidExtension) {
      setStatus(
        "Unsupported file format. Use PDF, JPG, JPEG, PNG, or WEBP receipt/proof documents.",
      );
      return;
    }

    const formData = new FormData();
    formData.append("proof", file);
    formData.append("donationId", selectedDonation);

    try {
      setIsUploading(true);
      setStatus("Uploading...");
      const res = await api.post("/ngo/proof", formData, {
        headers: {
          "X-ROLE": "NGO",
          "X-WALLET": walletAddress,
          "Content-Type": "multipart/form-data",
        },
      });

      const ipfs = res.data.ipfsHash as string;
      const msg = res.data.message as string | undefined;
      const extracted = res.data.extractedText as string | undefined;
      const note = res.data.extractionNote as string | undefined;
      let detail = `Uploaded successfully. IPFS: ${ipfs}\n\n${msg ?? "Proof saved."}`;
      if (extracted) {
        detail += `\n\nExtracted text (preview):\n${extracted.slice(0, 1500)}${extracted.length > 1500 ? "…" : ""}`;
      } else if (note) {
        detail += `\n\nExtraction: ${note}`;
      }
      detail +=
        "\n\nAudit status: pending admin review — an administrator must run the audit from the admin dashboard.";
      setStatus(detail);
    } catch (err) {
      console.error(err);
      const errorMessage =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message ||
        "Upload failed. Please check the proof format and backend configuration.";
      setStatus(errorMessage);
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">NGO Dashboard</h1>

      <div className="mb-4 flex gap-2">
        <input
          type="text"
          value={walletInput}
          onChange={(e) => setWalletInput(e.target.value)}
          placeholder="Enter NGO MetaMask wallet address (0x...)"
          className="border p-2 rounded w-full max-w-2xl"
        />
        <button
          onClick={loadTransactions}
          className="bg-black text-white px-4 py-2 rounded"
        >
          View Transactions
        </button>
      </div>
      {walletAddress && (
        <p className="mb-3">
          Showing transactions for NGO wallet: <b>{walletAddress}</b>
        </p>
      )}
      {transactionsError && <p className="text-red-600 mb-3">{transactionsError}</p>}

      {walletAddress && (
        <>
          {loadingTransactions ? (
            <p className="mb-4">Loading transactions...</p>
          ) : donations.length === 0 ? (
            <p className="mb-4">No transactions found for this NGO account yet.</p>
          ) : (
            <table className="w-full border mb-4">
              <thead>
                <tr className="bg-gray-100">
                  <th className="border p-2">Transaction Hash</th>
                  <th className="border p-2">Donor</th>
                  <th className="border p-2">Purpose</th>
                  <th className="border p-2">Amount</th>
                  <th className="border p-2">Status</th>
                </tr>
              </thead>
              <tbody>
                {donations.map((d) => (
                  <tr key={d.id}>
                    <td className="border p-2">{d.tx_hash}</td>
                    <td className="border p-2">{d.donor_wallet}</td>
                    <td className="border p-2">{d.purpose}</td>
                    <td className="border p-2">{d.amount}</td>
                    <td className="border p-2">{d.status}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          <select
            className="border p-2 mb-4 w-full"
            value={selectedDonation}
            onChange={(e) => setSelectedDonation(e.target.value)}
          >
            <option value="">Select Donation</option>
            {donations.map((d) => (
              <option key={d.id} value={d.id}>
                {d.purpose} — ₹{d.amount}
              </option>
            ))}
          </select>

          <div className="mb-4 rounded border border-blue-200 bg-blue-50 p-3 text-sm text-blue-900">
            <p className="font-semibold">Supported proof documents</p>
            <p>
              Upload only receipt/proof files in: <b>PDF, JPG, JPEG, PNG, WEBP</b>.
              Random/unrelated documents may be flagged during audit.
            </p>
          </div>

          <div className="mb-4 rounded border border-gray-300 p-4">
            <p className="mb-2 text-sm font-medium text-gray-700">Choose receipt/proof file</p>
            <div className="flex flex-wrap items-center gap-3">
              <input
                ref={fileInputRef}
                type="file"
                accept=".pdf,.jpg,.jpeg,.png,.webp,application/pdf,image/jpeg,image/png,image/webp"
                className="hidden"
                onChange={(e) => setFile(e.target.files?.[0] || null)}
              />
              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                className="rounded border border-gray-400 bg-white px-4 py-2 font-medium hover:bg-gray-50"
              >
                Choose File
              </button>
              <span className="text-sm text-gray-700">
                {file ? file.name : "No file selected"}
              </span>
            </div>
          </div>

          <button
            onClick={uploadProof}
            disabled={isUploading}
            className="rounded bg-blue-600 px-5 py-2 font-semibold text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-blue-300"
          >
            {isUploading ? "Uploading proof..." : "Upload Receipt / Proof"}
          </button>

          {status && (
            <p className="mt-4 whitespace-pre-wrap text-sm">{status}</p>
          )}
        </>
      )}
    </div>
  );
}
