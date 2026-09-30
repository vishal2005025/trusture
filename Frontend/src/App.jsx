import React, { useState } from 'react';
import WalletConnect from './components/WalletConnect';
import DonorDashboard from './components/DonorDashboard';
import NGODashboard from './components/NGODashboard';
import './App.css';

function App() {
  const [currentView, setCurrentView] = useState('home');
  const [isWalletConnected, setIsWalletConnected] = useState(false);

  const renderCurrentView = () => {
    switch (currentView) {
      case 'wallet':
        return <WalletConnect />;
      case 'donor':
        return <DonorDashboard />;
      case 'ngo':
        return <NGODashboard />;
      default:
        return <HomePage />;
    }
  };

  const HomePage = () => (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-secondary-50">
      {/* Navigation Header */}
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex justify-between items-center">
            <div className="flex items-center">
              <div className="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center mr-3">
                <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.818-4.818A4 4 0 0121 12a4 4 0 01-4 4H3" />
                </svg>
              </div>
              <h1 className="text-2xl font-bold text-gray-900">TrustuRe</h1>
              <span className="ml-3 text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded">
                Blockchain NGO Auditing
              </span>
            </div>
            <div className="flex space-x-4">
              <button
                onClick={() => setCurrentView('wallet')}
                className="btn-primary"
              >
                Connect Wallet
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <div className="relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-6 py-24">
          <div className="text-center">
            <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
              Transparent
              <span className="text-primary-600"> Blockchain</span>
              <br />
              NGO Donations
            </h1>
            <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto leading-relaxed">
              Revolutionizing charitable giving through blockchain technology. 
              Track every donation, ensure transparency, and build trust between donors and NGOs.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button
                onClick={() => setCurrentView('donor')}
                className="btn-primary text-lg px-8 py-4"
              >
                Start Donating
              </button>
              <button
                onClick={() => setCurrentView('ngo')}
                className="bg-white border-2 border-gray-300 hover:border-gray-400 text-gray-700 font-medium text-lg px-8 py-4 rounded-lg transition-colors duration-200"
              >
                NGO Dashboard
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="bg-white py-24">
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Why Choose TrustuRe?
            </h2>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">
              Built on blockchain technology to ensure complete transparency and trust in charitable donations.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="card text-center">
              <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <svg className="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.818-4.818A4 4 0 0121 12a4 4 0 01-4 4H3" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">Blockchain Transparency</h3>
              <p className="text-gray-600">
                Every transaction is recorded on the blockchain, providing immutable proof of where your donations go.
              </p>
            </div>

            <div className="card text-center">
              <div className="w-16 h-16 bg-secondary-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <svg className="w-8 h-8 text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">Real-time Analytics</h3>
              <p className="text-gray-600">
                Track donation impact with detailed analytics and reports powered by our Go backend.
              </p>
            </div>

            <div className="card text-center">
              <div className="w-16 h-16 bg-yellow-100 rounded-full flex items-center justify-center mx-auto mb-6">
                <svg className="w-8 h-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-4">Verified NGOs</h3>
              <p className="text-gray-600">
                All NGOs are verified and audited to ensure legitimacy and proper fund utilization.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* How It Works Section */}
      <div className="bg-gray-50 py-24">
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              How It Works
            </h2>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">
              Simple steps to make a transparent, traceable donation.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                1
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Connect Wallet</h3>
              <p className="text-gray-600">Connect your crypto wallet (MetaMask, WalletConnect)</p>
            </div>

            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                2
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Choose NGO</h3>
              <p className="text-gray-600">Browse and select from verified NGOs and campaigns</p>
            </div>

            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                3
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Make Donation</h3>
              <p className="text-gray-600">Send cryptocurrency donation through smart contracts</p>
            </div>

            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                4
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Track Impact</h3>
              <p className="text-gray-600">Monitor how your donation is used with blockchain transparency</p>
            </div>
          </div>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-gray-900 text-white py-12">
        <div className="max-w-7xl mx-auto px-6">
          <div className="text-center">
            <div className="flex items-center justify-center mb-4">
              <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center mr-2">
                <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.818-4.818A4 4 0 0121 12a4 4 0 01-4 4H3" />
                </svg>
              </div>
              <h3 className="text-xl font-bold">TrustuRe</h3>
            </div>
            <p className="text-gray-400 mb-4">
              Blockchain-Based NGO Donation Auditing Framework
            </p>
            <p className="text-sm text-gray-500">
              Built with React, Tailwind CSS, and ready for Go backend integration
            </p>
          </div>
        </div>
      </footer>
    </div>
  );

  return (
    <div className="App">
      {currentView !== 'home' && (
        <nav className="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-50">
          <div className="max-w-7xl mx-auto px-6 py-4">
            <div className="flex justify-between items-center">
              <button
                onClick={() => setCurrentView('home')}
                className="flex items-center text-primary-600 hover:text-primary-800"
              >
                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                </svg>
                Back to Home
              </button>
              <div className="flex items-center">
                <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center mr-2">
                  <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.818-4.818A4 4 0 0121 12a4 4 0 01-4 4H3" />
                  </svg>
                </div>
                <h1 className="text-xl font-bold text-gray-900">TrustuRe</h1>
              </div>
            </div>
          </div>
        </nav>
      )}
      
      {renderCurrentView()}
    </div>
  );
}

export default App;
