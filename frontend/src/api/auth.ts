import api from "./axios";

export async function getNonce(wallet: string) {
  const res = await api.get(`/auth/nonce?wallet=${wallet}`);
  return res.data.nonce;
}

export async function verifySignature(
  wallet: string,
  signature: string
) {
  const res = await api.get(
    `/auth/verify?wallet=${wallet}&signature=${signature}`
  );
  return res.data;
}
