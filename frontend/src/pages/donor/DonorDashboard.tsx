import { useEffect, useState } from "react";
import api from "../../api/axios";
import { useWallet } from "../../auth/useWallet";

type Donation = {
  id: string;
  amount: number;
  purpose: string;
  status: string;
  riskLevel: string;
};

export default function DonorDashboard() {
  const { address, isVerified, connectAndLogin } = useWallet();
  const [donations, setDonations] = useState<Donation[]>([]);

  useEffect(() => {
    if (!address || !isVerified) return;

    api
      .get("/donor/dashboard", {
        headers: { "X-ROLE": "DONOR", "X-WALLET": address },
      })
      .then((res) => setDonations(res.data))
      .catch(console.error);
  }, [address, isVerified]);

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Donor Dashboard</h1>

      {!address && (
        <button
          onClick={connectAndLogin}
          className="bg-black text-white px-4 py-2 rounded"
        >
          Connect & Login
        </button>
      )}

      {address && (
        <>
          <p className="mb-4">
            Connected wallet: <b>{address}</b>
          </p>

          <table className="w-full border">
            <thead>
              <tr className="bg-gray-100">
                <th className="border p-2">Purpose</th>
                <th className="border p-2">Amount</th>
                <th className="border p-2">Status</th>
                <th className="border p-2">Risk</th>
              </tr>
            </thead>
            <tbody>
              {donations.map((d) => (
                <tr key={d.id}>
                  <td className="border p-2">{d.purpose}</td>
                  <td className="border p-2">{d.amount}</td>
                  <td className="border p-2">{d.status}</td>
                  <td className="border p-2">{d.riskLevel}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  );
}
