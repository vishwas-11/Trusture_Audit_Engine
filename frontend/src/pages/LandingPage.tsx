// src/pages/LandingPage.tsx
// import React from "react";
import { useNavigate } from "react-router-dom";
import BlockchainScene from "../components/BlockchainScene";
// import { ArrowRight, ShieldCheck, Globe, Users } from "lucide-react"; // Assuming you might have icons, if not, standard text works nicely too.

const LandingPage = () => {
  const navigate = useNavigate();

  return (
    <div className="relative w-full h-screen overflow-hidden text-white font-sans">
      {/* 3D Background */}
      <BlockchainScene />

      {/* Main Content Overlay */}
      <div className="relative z-10 flex flex-col h-full overflow-y-auto">
        
        {/* Navbar */}
        <nav className="flex justify-between items-center px-8 py-6 backdrop-blur-sm">
          <div className="text-2xl font-bold tracking-tighter text-sky-400">
            TRUSTURE
          </div>
          <div className="space-x-6 text-sm font-medium text-slate-300">
            <button className="hover:text-white transition-colors">About</button>
            <button className="hover:text-white transition-colors">Audit</button>
            <button className="px-4 py-2 border border-slate-600 rounded-full hover:bg-sky-500 hover:border-sky-500 hover:text-white transition-all duration-300">
              Connect Wallet
            </button>
          </div>
        </nav>

        {/* Hero Section */}
        <div className="flex-grow flex flex-col items-center justify-center text-center px-4">
          <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight mb-6 bg-clip-text text-transparent bg-gradient-to-r from-sky-400 to-indigo-500">
            Trust. Trace. Transform.
          </h1>
          <p className="max-w-2xl text-lg md:text-xl text-slate-400 mb-10 leading-relaxed">
            The decentralized ecosystem ensuring 100% transparency for NGOs and Donors.
            Track every transaction on the immutable blockchain.
          </p>
        </div>

        {/* Interactive Cards Section */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-6xl mx-auto w-full px-6 pb-20">
          
          {/* Card 1: Donor */}
          <div 
            onClick={() => navigate('/donor')}
            className="group relative p-8 rounded-2xl bg-slate-900/40 border border-slate-700/50 backdrop-blur-md cursor-pointer hover:-translate-y-2 transition-transform duration-300 hover:border-sky-500/50 hover:shadow-lg hover:shadow-sky-500/20"
          >
            <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
              {/* Abstract Icon Placeholder */}
              <div className="w-24 h-24 bg-sky-400 rounded-full blur-2xl"></div>
            </div>
            <h3 className="text-2xl font-bold text-white mb-2">I am a Donor</h3>
            <p className="text-slate-400 mb-6">Make secure donations and track the impact of your funds in real-time.</p>
            <div className="flex items-center text-sky-400 font-semibold group-hover:gap-2 transition-all">
              Start Giving <span>&rarr;</span>
            </div>
          </div>

          {/* Card 2: NGO */}
          <div 
            onClick={() => navigate('/ngo')}
            className="group relative p-8 rounded-2xl bg-slate-900/40 border border-slate-700/50 backdrop-blur-md cursor-pointer hover:-translate-y-2 transition-transform duration-300 hover:border-purple-500/50 hover:shadow-lg hover:shadow-purple-500/20"
          >
             <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
              <div className="w-24 h-24 bg-purple-400 rounded-full blur-2xl"></div>
            </div>
            <h3 className="text-2xl font-bold text-white mb-2">I am an NGO</h3>
            <p className="text-slate-400 mb-6">Register your organization, create campaigns, and validate milestones.</p>
            <div className="flex items-center text-purple-400 font-semibold group-hover:gap-2 transition-all">
              Join Network <span>&rarr;</span>
            </div>
          </div>

          {/* Card 3: Admin */}
          <div 
            onClick={() => navigate('/admin')}
            className="group relative p-8 rounded-2xl bg-slate-900/40 border border-slate-700/50 backdrop-blur-md cursor-pointer hover:-translate-y-2 transition-transform duration-300 hover:border-emerald-500/50 hover:shadow-lg hover:shadow-emerald-500/20"
          >
             <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
              <div className="w-24 h-24 bg-emerald-400 rounded-full blur-2xl"></div>
            </div>
            <h3 className="text-2xl font-bold text-white mb-2">Auditor / Admin</h3>
            <p className="text-slate-400 mb-6">Oversee the ecosystem, verify NGO credentials, and ensure compliance.</p>
            <div className="flex items-center text-emerald-400 font-semibold group-hover:gap-2 transition-all">
              Access Portal <span>&rarr;</span>
            </div>
          </div>

        </div>
      </div>
    </div>
  );
};

export default LandingPage;