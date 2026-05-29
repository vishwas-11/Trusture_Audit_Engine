import { useEffect, useState } from "react";
import api from "../../api/axios";

const TOKEN_KEY = "trusture_admin_token";

type OverviewProof = {
  id: string;
  ipfs_hash: string;
  type: string;
  extracted_text: string;
  extraction_note: string;
  uploaded_at: number;
};

type OverviewAudit = {
  id?: string;
  donation_id?: string;
  status: string;
  score: number;
  flags: string[];
  created_at: number;
  donor_ngo_distance_km?: number;
  ngo_proof_distance_km?: number;
  gst_verification_source?: string;
};

type AuditResponse = {
  id?: string;
  donation_id?: string;
  status?: string;
  score?: number;
  flags?: string[];
  created_at?: number;
  donor_ngo_distance_km?: number;
  ngo_proof_distance_km?: number;
  gst_verification_source?: string;
};

type OverviewRow = {
  id: string;
  tx_hash: string;
  donor_wallet: string;
  ngo_wallet: string;
  amount: number;
  purpose: string;
  status: string;
  created_at: number;
  proofs: OverviewProof[];
  latest_audit?: OverviewAudit;
};

type FilterMode = "all" | "audited" | "not_audited" | "with_proof" | "without_proof";

function normalizeRows(data: unknown): OverviewRow[] {
  if (!Array.isArray(data)) return [];

  return data.map((row) => {
    const record = (row ?? {}) as Partial<OverviewRow> & { proofs?: OverviewProof[] | null };
    return {
      id: record.id ?? "",
      tx_hash: record.tx_hash ?? "",
      donor_wallet: record.donor_wallet ?? "",
      ngo_wallet: record.ngo_wallet ?? "",
      amount: record.amount ?? 0,
      purpose: record.purpose ?? "",
      status: record.status ?? "",
      created_at: record.created_at ?? 0,
      proofs: Array.isArray(record.proofs) ? record.proofs : [],
      latest_audit: record.latest_audit,
    };
  });
}

function normalizeAudit(data: unknown): OverviewAudit | undefined {
  const record = (data ?? {}) as Partial<AuditResponse>;
  if (!record.status) return undefined;

  return {
    id: record.id ?? "",
    donation_id: record.donation_id ?? "",
    status: record.status,
    score: record.score ?? 0,
    flags: Array.isArray(record.flags) ? record.flags : [],
    created_at: record.created_at ?? 0,
    donor_ngo_distance_km: record.donor_ngo_distance_km,
    ngo_proof_distance_km: record.ngo_proof_distance_km,
    gst_verification_source: record.gst_verification_source,
  };
}

function normalizeAuditHistory(data: unknown): OverviewAudit[] {
  if (!Array.isArray(data)) return [];

  return data
    .map((item) => normalizeAudit(item))
    .filter((item): item is OverviewAudit => Boolean(item));
}

export default function AdminDashboard() {
  const [token, setToken] = useState<string>(() => localStorage.getItem(TOKEN_KEY) || "");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [mode, setMode] = useState<"login" | "register">("login");
  const [authError, setAuthError] = useState("");
  const [loadingAuth, setLoadingAuth] = useState(false);

  const [rows, setRows] = useState<OverviewRow[]>([]);
  const [overviewError, setOverviewError] = useState("");
  const [loadingOverview, setLoadingOverview] = useState(false);
  const [auditBusyId, setAuditBusyId] = useState<string | null>(null);
  const [filterMode, setFilterMode] = useState<FilterMode>("all");
  const [expandedHistoryId, setExpandedHistoryId] = useState<string | null>(null);
  const [historyByDonation, setHistoryByDonation] = useState<Record<string, OverviewAudit[]>>({});
  const [historyLoadingId, setHistoryLoadingId] = useState<string | null>(null);

  useEffect(() => {
    if (!token) return;
    setLoadingOverview(true);
    setOverviewError("");
    api
      .get("/admin/overview", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((res) => setRows(normalizeRows(res.data)))
      .catch((err) => {
        console.error(err);
        setOverviewError("Failed to load overview. Check token or server.");
        setRows([]);
      })
      .finally(() => setLoadingOverview(false));
  }, [token]);

  const persistToken = (t: string) => {
    localStorage.setItem(TOKEN_KEY, t);
    setToken(t);
  };

  const logout = () => {
    localStorage.removeItem(TOKEN_KEY);
    setToken("");
    setRows([]);
    setHistoryByDonation({});
    setExpandedHistoryId(null);
  };

  const submitAuth = async () => {
    setAuthError("");
    setLoadingAuth(true);
    try {
      const path = mode === "login" ? "/admin/auth/login" : "/admin/auth/register";
      const res = await api.post(path, { email: email.trim(), password });
      if (mode === "login") {
        persistToken(res.data.token);
      } else {
        setMode("login");
        setAuthError("Account created. You can sign in now.");
      }
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message ||
        "Request failed";
      setAuthError(msg);
    } finally {
      setLoadingAuth(false);
    }
  };

  const runAudit = async (donationId: string) => {
    if (!token) return;
    setAuditBusyId(donationId);
    setOverviewError("");
    try {
      const res = await api.post(
        `/admin/donations/${donationId}/audit`,
        {},
        { headers: { Authorization: `Bearer ${token}` } },
      );
      const audit = normalizeAudit(res.data?.audit);
      if (audit) {
        setRows((current) =>
          current.map((row) => (row.id === donationId ? { ...row, latest_audit: audit } : row)),
        );
        setHistoryByDonation((current) => ({
          ...current,
          [donationId]: [audit, ...(current[donationId] ?? [])],
        }));
      }

      const overview = await api.get("/admin/overview", {
        headers: { Authorization: `Bearer ${token}` },
      });
      setRows(normalizeRows(overview.data));
    } catch (err: unknown) {
      console.error(err);
      const msg =
        (err as { response?: { data?: { message?: string } } })?.response?.data?.message ||
        "Audit failed";
      setOverviewError(msg);
    } finally {
      setAuditBusyId(null);
    }
  };

  const loadAuditHistory = async (donationId: string) => {
    if (!token || historyByDonation[donationId]) return;

    setHistoryLoadingId(donationId);
    try {
      const res = await api.get(`/admin/donations/${donationId}/audits`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setHistoryByDonation((current) => ({
        ...current,
        [donationId]: normalizeAuditHistory(res.data),
      }));
    } catch (err) {
      console.error(err);
      setOverviewError("Failed to load audit history for that transaction.");
    } finally {
      setHistoryLoadingId(null);
    }
  };

  const toggleHistory = async (donationId: string) => {
    setOverviewError("");
    setExpandedHistoryId((current) => (current === donationId ? null : donationId));
    await loadAuditHistory(donationId);
  };

  const displayedRows = rows.filter((row) => {
    switch (filterMode) {
      case "audited":
        return Boolean(row.latest_audit);
      case "not_audited":
        return !row.latest_audit;
      case "with_proof":
        return row.proofs.length > 0;
      case "without_proof":
        return row.proofs.length === 0;
      default:
        return true;
    }
  });

  if (!token) {
    return (
      <div className="mx-auto max-w-md p-6">
        <h1 className="mb-4 text-2xl font-bold">Admin dashboard</h1>
        <p className="mb-4 text-sm text-gray-600">
          Sign in with your admin email and password. Audits run only after you trigger them here — proof
          uploads no longer auto-audit.
        </p>
        <div className="mb-4 flex gap-2">
          <button
            type="button"
            className={`rounded px-3 py-1 ${mode === "login" ? "bg-black text-white" : "bg-gray-200"}`}
            onClick={() => setMode("login")}
          >
            Sign in
          </button>
          <button
            type="button"
            className={`rounded px-3 py-1 ${mode === "register" ? "bg-black text-white" : "bg-gray-200"}`}
            onClick={() => setMode("register")}
          >
            Register
          </button>
        </div>
        <label className="mb-2 block text-sm font-medium">Email</label>
        <input
          className="mb-3 w-full rounded border p-2"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          autoComplete="email"
        />
        <label className="mb-2 block text-sm font-medium">Password</label>
        <input
          className="mb-4 w-full rounded border p-2"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          autoComplete={mode === "login" ? "current-password" : "new-password"}
        />
        {authError && <p className="mb-3 text-sm text-red-600">{authError}</p>}
        <button
          type="button"
          disabled={loadingAuth}
          className="w-full rounded bg-black py-2 font-semibold text-white disabled:opacity-50"
          onClick={submitAuth}
        >
          {loadingAuth ? "Please wait…" : mode === "login" ? "Sign in" : "Create admin account"}
        </button>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6 flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-2xl font-bold">Admin — transactions & audits</h1>
        <button type="button" className="rounded border border-gray-400 px-4 py-2 text-sm" onClick={logout}>
          Sign out
        </button>
      </div>

      <p className="mb-4 max-w-3xl text-sm text-gray-700">
        Each row shows blockchain-recorded donations, extracted proof text (when available), and the latest audit
        once you run it. Use &quot;Run audit&quot; after NGOs upload proofs.
      </p>

      <div className="mb-4 flex flex-wrap items-center gap-3 rounded-lg border border-slate-200 bg-white p-3 shadow-sm">
        <label className="text-sm font-medium text-slate-700" htmlFor="audit-filter">
          View
        </label>
        <select
          id="audit-filter"
          value={filterMode}
          onChange={(e) => setFilterMode(e.target.value as FilterMode)}
          className="rounded border border-slate-300 bg-white px-3 py-2 text-sm"
        >
          <option value="all">All transactions</option>
          <option value="audited">Audited transactions</option>
          <option value="not_audited">Non-audited transactions</option>
          <option value="with_proof">Transactions with proof uploaded</option>
          <option value="without_proof">Transactions without proof</option>
        </select>
        <p className="text-xs text-slate-500">
          Showing {displayedRows.length} of {rows.length} records.
        </p>
      </div>

      {overviewError && <p className="mb-4 text-sm text-red-600">{overviewError}</p>}
      {loadingOverview ? (
        <p>Loading…</p>
      ) : displayedRows.length === 0 ? (
        <p>No donations recorded yet.</p>
      ) : (
        <div className="space-y-8">
          {displayedRows.map((row) => (
            <section key={row.id} className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
              <div className="mb-3 flex flex-wrap items-start justify-between gap-3">
                <div>
                  <p className="font-semibold text-gray-900">{row.purpose}</p>
                  <p className="text-xs text-gray-500">Donation ID: {row.id}</p>
                  <p className="text-xs text-gray-500 break-all">Tx: {row.tx_hash}</p>
                  <p className="mt-1 text-sm">
                    {row.amount} MATIC · Status: <span className="font-medium">{row.status}</span>
                  </p>
                  <p className="mt-1 text-sm text-gray-700">
                    <span className="font-medium">From:</span>{" "}
                    <span className="break-all font-mono text-xs">{row.donor_wallet}</span>
                  </p>
                  <p className="text-sm text-gray-700">
                    <span className="font-medium">To:</span>{" "}
                    <span className="break-all font-mono text-xs">{row.ngo_wallet}</span>
                  </p>
                </div>
                <button
                  type="button"
                  disabled={auditBusyId === row.id || row.proofs.length === 0}
                  title={row.proofs.length === 0 ? "No proof uploaded yet" : undefined}
                  className="shrink-0 rounded bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-blue-300"
                  onClick={() => runAudit(row.id)}
                >
                  {auditBusyId === row.id ? "Running…" : "Run audit"}
                </button>
              </div>

              <div className="mb-3 flex flex-wrap gap-2">
                <button
                  type="button"
                  className="rounded border border-slate-300 px-3 py-1.5 text-xs font-semibold text-slate-700 hover:bg-slate-50"
                  onClick={() => toggleHistory(row.id)}
                >
                  {expandedHistoryId === row.id ? "Hide audit history" : "Show audit history"}
                </button>
                {historyLoadingId === row.id ? (
                  <span className="rounded bg-slate-100 px-3 py-1.5 text-xs text-slate-600">
                    Loading history...
                  </span>
                ) : null}
              </div>

              {auditBusyId === row.id ? (
                <div className="mb-3 rounded border border-blue-200 bg-blue-100 p-3 text-sm text-blue-900">
                  Audit is running. We’ll refresh this row as soon as the result is stored.
                </div>
              ) : null}

              <div className="mb-3 rounded border border-blue-100 bg-blue-50 p-3 text-sm text-blue-900">
                Uploaded proofs for this transaction: <span className="font-semibold">{row.proofs.length}</span>
              </div>

              <div className="mb-3 rounded border border-slate-100 bg-slate-50 p-3">
                <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-slate-600">Latest audit</p>
                {row.latest_audit ? (
                  <div className="text-sm">
                    <p>
                      <span className="text-gray-600">Result:</span>{" "}
                      <span className="font-semibold">{row.latest_audit.status}</span> · Score:{" "}
                      {row.latest_audit.score?.toFixed?.(2) ?? row.latest_audit.score}
                    </p>
                    <p className="mt-1 text-xs text-gray-600">
                      Last run:{" "}
                      {row.latest_audit.created_at
                        ? new Date(row.latest_audit.created_at * 1000).toLocaleString()
                        : "Unknown"}
                    </p>
                    {row.latest_audit.flags?.length ? (
                      <p className="mt-1 text-xs text-gray-700">Flags: {row.latest_audit.flags.join(", ")}</p>
                    ) : null}
                  </div>
                ) : (
                  <p className="text-sm text-amber-800">No audit run yet — click &quot;Run audit&quot; when ready.</p>
                )}
              </div>

              {expandedHistoryId === row.id ? (
                <div className="mb-3 rounded border border-emerald-100 bg-emerald-50 p-3">
                  <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-emerald-700">
                    Audit history
                  </p>
                  {historyByDonation[row.id]?.length ? (
                    <div className="space-y-3">
                      {historyByDonation[row.id].map((audit, index) => (
                        <div
                          key={`${audit.id ?? row.id}-${audit.created_at}-${index}`}
                          className="rounded border border-emerald-100 bg-white p-3 text-sm"
                        >
                          <div className="flex flex-wrap items-center justify-between gap-2">
                            <p className="font-semibold text-slate-900">
                              {audit.status} · Score {audit.score.toFixed(2)}
                            </p>
                            <p className="text-xs text-slate-500">
                              {audit.created_at
                                ? new Date(audit.created_at * 1000).toLocaleString()
                                : "Unknown"}
                            </p>
                          </div>
                          {audit.flags.length ? (
                            <p className="mt-1 text-xs text-slate-700">Flags: {audit.flags.join(", ")}</p>
                          ) : (
                            <p className="mt-1 text-xs text-slate-600">No flags recorded.</p>
                          )}
                          <div className="mt-2 grid gap-2 text-xs text-slate-600 sm:grid-cols-2">
                            <p>
                              Donor-NGO distance:{" "}
                              {typeof audit.donor_ngo_distance_km === "number"
                                ? `${audit.donor_ngo_distance_km.toFixed(2)} km`
                                : "N/A"}
                            </p>
                            <p>
                              NGO-proof distance:{" "}
                              {typeof audit.ngo_proof_distance_km === "number"
                                ? `${audit.ngo_proof_distance_km.toFixed(2)} km`
                                : "N/A"}
                            </p>
                            <p>GST source: {audit.gst_verification_source || "N/A"}</p>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : historyLoadingId === row.id ? (
                    <p className="text-sm text-emerald-800">Loading audit history...</p>
                  ) : (
                    <p className="text-sm text-emerald-800">No audit history recorded yet.</p>
                  )}
                </div>
              ) : null}

              <div>
                <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-600">
                  Proof extraction (for admin review)
                </p>
                {row.proofs.length === 0 ? (
                  <p className="text-sm text-gray-600">No proofs uploaded for this donation.</p>
                ) : (
                  row.proofs.map((p) => (
                    <div key={p.id} className="mb-3 rounded border border-gray-200 p-3 last:mb-0">
                      <p className="text-xs text-gray-500">
                        Proof · {p.type} · IPFS {p.ipfs_hash.slice(0, 12)}…
                      </p>
                      {p.extracted_text ? (
                        <pre className="mt-2 max-h-48 overflow-auto whitespace-pre-wrap rounded bg-gray-50 p-2 text-xs text-gray-900">
                          {p.extracted_text}
                        </pre>
                      ) : (
                        <p className="mt-2 text-sm text-amber-900">
                          {p.extraction_note ||
                            "No relevant information could be extracted from this document."}
                        </p>
                      )}
                    </div>
                  ))
                )}
              </div>
            </section>
          ))}
        </div>
      )}
    </div>
  );
}
