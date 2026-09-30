import React, { useState } from 'react';

const WalletConnect = () => {
  const [isConnected, setIsConnected] = useState(false);
  const [walletAddress, setWalletAddress] = useState('');
  const [balance, setBalance] = useState('0.00');
  const [isConnecting, setIsConnecting] = useState(false);

  // Placeholder function for wallet connection
  const connectWallet = async () => {
    setIsConnecting(true);
    try {
      // TODO: Integrate with Web3 wallet (MetaMask, WalletConnect, etc.)
      // This is a placeholder that simulates wallet connection
      setTimeout(() => {
        setIsConnected(true);
        setWalletAddress('0x1234...5678'); // Mock address
        setBalance('2.45'); // Mock balance
        setIsConnecting(false);
      }, 2000);
    } catch (error) {
      console.error('Failed to connect wallet:', error);
      setIsConnecting(false);
    }
  };

  const disconnectWallet = () => {
    setIsConnected(false);
    setWalletAddress('');
    setBalance('0.00');
  };

  return (
    <div className="card max-w-md mx-auto">
      <div className="text-center">
        <div className="mb-6">
          <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Wallet Connection</h2>
          <p className="text-gray-600">
            Connect your crypto wallet to start donating and tracking contributions
          </p>
        </div>

        {!isConnected ? (
          <div className="space-y-4">
            <button
              onClick={connectWallet}
              disabled={isConnecting}
              className={`w-full btn-primary ${isConnecting ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              {isConnecting ? (
                <div className="flex items-center justify-center">
                  <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Connecting...
                </div>
              ) : (
                'Connect Wallet'
              )}
            </button>
            
            <div className="text-sm text-gray-500">
              <p className="mb-2">Supported wallets:</p>
              <div className="flex justify-center space-x-4">
                <div className="flex items-center space-x-1">
                  <div className="w-6 h-6 bg-orange-500 rounded"></div>
                  <span>MetaMask</span>
                </div>
                <div className="flex items-center space-x-1">
                  <div className="w-6 h-6 bg-blue-500 rounded"></div>
                  <span>WalletConnect</span>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="bg-green-50 border border-green-200 rounded-lg p-4">
              <div className="flex items-center justify-center mb-2">
                <svg className="w-5 h-5 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
                </svg>
                <span className="text-green-800 font-medium">Wallet Connected</span>
              </div>
              
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Address:</span>
                  <span className="font-mono text-gray-900">{walletAddress}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Balance:</span>
                  <span className="font-bold text-gray-900">{balance} ETH</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Network:</span>
                  <span className="text-gray-900">Ethereum Mainnet</span>
                </div>
              </div>
            </div>

            <button
              onClick={disconnectWallet}
              className="w-full bg-gray-500 hover:bg-gray-600 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200"
            >
              Disconnect Wallet
            </button>
          </div>
        )}

        <div className="mt-6 p-4 bg-blue-50 rounded-lg">
          <div className="flex items-start">
            <svg className="w-5 h-5 text-blue-500 mr-2 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div className="text-left">
              <p className="text-sm font-medium text-blue-800 mb-1">Blockchain Integration</p>
              <p className="text-xs text-blue-600">
                Your donations will be recorded on the blockchain for complete transparency and auditing.
                Smart contracts ensure funds reach their intended recipients.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default WalletConnect;