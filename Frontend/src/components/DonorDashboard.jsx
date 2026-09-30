import React, { useState } from 'react';

const DonorDashboard = () => {
  const [activeTab, setActiveTab] = useState('discover');
  
  // Mock data - will be replaced with API calls to Go backend
  const [donationHistory] = useState([
    {
      id: 1,
      ngoName: 'Save the Children',
      amount: '0.5 ETH',
      usdAmount: '$1,250',
      date: '2025-09-28',
      status: 'Completed',
      txHash: '0xabc123...def456',
      category: 'Education'
    },
    {
      id: 2,
      ngoName: 'Red Cross',
      amount: '1.0 ETH',
      usdAmount: '$2,500',
      date: '2025-09-25',
      status: 'Completed',
      txHash: '0x789xyz...abc123',
      category: 'Healthcare'
    },
    {
      id: 3,
      ngoName: 'World Wildlife Fund',
      amount: '0.3 ETH',
      usdAmount: '$750',
      date: '2025-09-20',
      status: 'Pending',
      txHash: '0xdef456...789xyz',
      category: 'Environment'
    }
  ]);

  const [availableNGOs] = useState([
    {
      id: 1,
      name: 'Save the Children',
      description: 'Protecting children worldwide through education and healthcare initiatives.',
      category: 'Education',
      verified: true,
      totalRaised: '245.8 ETH',
      donorsCount: 1247,
      transparency: 98,
      image: '🎓'
    },
    {
      id: 2,
      name: 'Red Cross',
      description: 'Emergency relief and disaster response organization helping communities in crisis.',
      category: 'Healthcare',
      verified: true,
      totalRaised: '567.2 ETH',
      donorsCount: 2156,
      transparency: 95,
      image: '🏥'
    },
    {
      id: 3,
      name: 'World Wildlife Fund',
      description: 'Conservation efforts to protect endangered species and their habitats.',
      category: 'Environment',
      verified: true,
      totalRaised: '189.4 ETH',
      donorsCount: 892,
      transparency: 97,
      image: '🌍'
    },
    {
      id: 4,
      name: 'Doctors Without Borders',
      description: 'Medical humanitarian aid in conflict zones and underserved areas.',
      category: 'Healthcare',
      verified: true,
      totalRaised: '432.1 ETH',
      donorsCount: 1567,
      transparency: 96,
      image: '👨‍⚕️'
    }
  ]);

  const totalDonated = donationHistory.reduce((sum, donation) => {
    return sum + parseFloat(donation.amount.split(' ')[0]);
  }, 0);

  const handleDonate = (ngoId) => {
    // TODO: Integrate with smart contract for donation
    console.log(`Initiating donation to NGO ${ngoId}`);
    // This will trigger the donation flow with wallet integration
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Donor Dashboard</h1>
          <p className="text-gray-600">Track your donations and discover verified NGOs</p>
        </div>

        {/* Stats Overview */}
        <div className="dashboard-grid mb-8">
          <div className="card">
            <div className="flex items-center">
              <div className="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center mr-4">
                <svg className="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
                </svg>
              </div>
              <div>
                <p className="text-sm text-gray-600">Total Donated</p>
                <p className="text-2xl font-bold text-gray-900">{totalDonated.toFixed(2)} ETH</p>
                <p className="text-sm text-green-600">≈ ${(totalDonated * 2500).toLocaleString()}</p>
              </div>
            </div>
          </div>

          <div className="card">
            <div className="flex items-center">
              <div className="w-12 h-12 bg-secondary-100 rounded-lg flex items-center justify-center mr-4">
                <svg className="w-6 h-6 text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
              </div>
              <div>
                <p className="text-sm text-gray-600">NGOs Supported</p>
                <p className="text-2xl font-bold text-gray-900">{new Set(donationHistory.map(d => d.ngoName)).size}</p>
                <p className="text-sm text-blue-600">Verified organizations</p>
              </div>
            </div>
          </div>

          <div className="card">
            <div className="flex items-center">
              <div className="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center mr-4">
                <svg className="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <p className="text-sm text-gray-600">Impact Score</p>
                <p className="text-2xl font-bold text-gray-900">97%</p>
                <p className="text-sm text-yellow-600">Transparency rating</p>
              </div>
            </div>
          </div>
        </div>

        {/* Navigation Tabs */}
        <div className="border-b border-gray-200 mb-6">
          <nav className="-mb-px flex space-x-8">
            <button
              onClick={() => setActiveTab('discover')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'discover'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Discover NGOs
            </button>
            <button
              onClick={() => setActiveTab('history')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'history'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Donation History
            </button>
          </nav>
        </div>

        {/* Content based on active tab */}
        {activeTab === 'discover' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold text-gray-900">Verified NGOs</h2>
              <div className="flex space-x-2">
                <select className="border border-gray-300 rounded-lg px-3 py-2 text-sm">
                  <option>All Categories</option>
                  <option>Education</option>
                  <option>Healthcare</option>
                  <option>Environment</option>
                  <option>Disaster Relief</option>
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {availableNGOs.map((ngo) => (
                <div key={ngo.id} className="card hover:shadow-xl transition-shadow duration-200">
                  <div className="flex items-start justify-between mb-4">
                    <div className="flex items-center">
                      <div className="text-3xl mr-3">{ngo.image}</div>
                      <div>
                        <h3 className="text-lg font-semibold text-gray-900 flex items-center">
                          {ngo.name}
                          {ngo.verified && (
                            <svg className="w-5 h-5 text-blue-500 ml-2" fill="currentColor" viewBox="0 0 20 20">
                              <path fillRule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                            </svg>
                          )}
                        </h3>
                        <span className="text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded">{ngo.category}</span>
                      </div>
                    </div>
                  </div>
                  
                  <p className="text-gray-600 mb-4">{ngo.description}</p>
                  
                  <div className="space-y-3 mb-4">
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Total Raised:</span>
                      <span className="font-semibold">{ngo.totalRaised}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Donors:</span>
                      <span className="font-semibold">{ngo.donorsCount.toLocaleString()}</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Transparency:</span>
                      <span className="font-semibold text-green-600">{ngo.transparency}%</span>
                    </div>
                  </div>
                  
                  <div className="flex space-x-3">
                    <button
                      onClick={() => handleDonate(ngo.id)}
                      className="flex-1 btn-primary"
                    >
                      Donate Now
                    </button>
                    <button className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors duration-200">
                      View Details
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'history' && (
          <div className="space-y-6">
            <h2 className="text-xl font-semibold text-gray-900">Your Donation History</h2>
            
            <div className="card overflow-hidden">
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">NGO</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Transaction</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {donationHistory.map((donation) => (
                      <tr key={donation.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div>
                            <div className="text-sm font-medium text-gray-900">{donation.ngoName}</div>
                            <div className="text-sm text-gray-500">{donation.category}</div>
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900">{donation.amount}</div>
                          <div className="text-sm text-gray-500">{donation.usdAmount}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {donation.date}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                            donation.status === 'Completed' 
                              ? 'bg-green-100 text-green-800' 
                              : 'bg-yellow-100 text-yellow-800'
                          }`}>
                            {donation.status}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <a 
                            href={`https://etherscan.io/tx/${donation.txHash}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-primary-600 hover:text-primary-900 font-mono"
                          >
                            {donation.txHash}
                          </a>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DonorDashboard;