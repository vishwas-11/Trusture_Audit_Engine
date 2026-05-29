import { BrowserRouter, Routes, Route } from "react-router-dom";
import DonorDashboard from "./pages/donor/DonorDashboard";
import NGODashboard from "./pages/ngo/NGODashboard";
import AdminDashboard from "./pages/admin/AdminDashboard";
import LandingPage from "./pages/LandingPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/donor" element={<DonorDashboard />} />
        <Route path="/ngo" element={<NGODashboard />} />
        <Route path="/admin" element={<AdminDashboard />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
