// import { useState } from "react";
// import { ethers } from "ethers";

// export function useWallet() {
//   const [address, setAddress] = useState<string | null>(null);

//   const connectWallet = async () => {
//     if (!(window as any).ethereum) {
//       alert("MetaMask not installed");
//       return;
//     }

//     const provider = new ethers.BrowserProvider(
//       (window as any).ethereum
//     );

//     const accounts = await provider.send("eth_requestAccounts", []);
//     setAddress(accounts[0]);
//   };

//   return { address, connectWallet };
// }



import { useState } from "react";
import { ethers } from "ethers";
import { getNonce, verifySignature } from "../api/auth";

export function useWallet() {
  const [address, setAddress] = useState<string | null>(null);
  const [isVerified, setIsVerified] = useState(false);

  const connectAndLogin = async () => {
    if (!(window as any).ethereum) {
      alert("MetaMask not installed");
      return;
    }

    const provider = new ethers.BrowserProvider(
      (window as any).ethereum
    );
    const signer = await provider.getSigner();
    const wallet = await signer.getAddress();
    setAddress(wallet);

    // Step 1: Get nonce from backend
    const nonce = await getNonce(wallet);

    // Step 2: Ask user to sign message
    const message = `Sign this message to login to TRUSTURE:\nNonce: ${nonce}\nWallet: ${wallet}`;
    const signature = await signer.signMessage(message);

    // Step 3: Verify signature with backend
    await verifySignature(wallet, signature);

    setIsVerified(true);
  };

  return { address, isVerified, connectAndLogin };
}
