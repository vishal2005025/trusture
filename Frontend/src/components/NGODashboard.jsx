import React, { useState, useEffect } from 'react';
import BlockchainViewer from './BlockchainViewer';
import TransactionTracker from './TransactionTracker';

const NGODashboard = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const [latestTransaction, setLatestTransaction] = useState(null);
  
  // Mock data - will be replaced with API calls to Go backend
  const [ngoProfile] = useState({
    name: 'Save the Children Foundation',
    description: 'Working to protect children worldwide through education, healthcare, and emergency relief programs.',
    category: 'Education & Child Welfare',
    established: '1919',
    location: 'London, UK',
    verified: true,
    transparencyScore: 98,
    taxId: 'EIN: 06-0726487'
  });

  const [campaigns] = useState([
    {
      id: 1,
      title: 'Emergency Education for Syrian Refugees',
      description: 'Providing educational resources and temporary schools for displaced children.',
      goal: '50 ETH',
      raised: '32.5 ETH',
      progress: 65,
      donors: 234,
      daysLeft: 15,
      status: 'Active',
      category: 'Education'
    },
    {
      id: 2,
      title: 'Clean Water Initiative - Africa',
      description: 'Building water wells and sanitation facilities in rural African communities.',
      goal: '75 ETH',
      raised: '68.2 ETH',
      progress: 91,
      donors: 456,
      daysLeft: 8,
      status: 'Active',
      category: 'Healthcare'
    },
    {
      id: 3,
      title: 'School Lunch Program - India',
      description: 'Providing nutritious meals to underprivileged children in Indian schools.',
      goal: '30 ETH',
      raised: '30 ETH',
      progress: 100,
      donors: 189,
      daysLeft: 0,
      status: 'Completed',
      category: 'Nutrition'
    }
  ]);

  const [recentDonations] = useState([
    {
      id: 1,
      donor: '0x1234...5678',
      amount: '2.5 ETH',
      usdAmount: '$6,250',
      campaign: 'Emergency Education for Syrian Refugees',
      date: '2025-09-30',
      txHash: '0xabc123...def456'
    },
    {
      id: 2,
      donor: '0x9876...4321',
      amount: '1.0 ETH',
      usdAmount: '$2,500',
      campaign: 'Clean Water Initiative - Africa',
      date: '2025-09-30',
      txHash: '0xdef456...abc123'
    },
    {
      id: 3,
      donor: '0x5555...7777',
      amount: '0.8 ETH',
      usdAmount: '$2,000',
      campaign: 'Emergency Education for Syrian Refugees',
      date: '2025-09-29',
      txHash: '0x789xyz...123abc'
    }
  ]);

  const totalRaised = campaigns.reduce((sum, campaign) => {
    return sum + parseFloat(campaign.raised.split(' ')[0]);
  }, 0);

  const totalDonors = campaigns.reduce((sum, campaign) => sum + campaign.donors, 0);

  const handleCreateCampaign = () => {
    // TODO: Integrate with Go backend to create new campaign
    console.log('Creating new campaign');
  };

  const handleWithdrawFunds = (campaignId) => {
    // TODO: Integrate with smart contract for fund withdrawal
    console.log(`Withdrawing funds from campaign ${campaignId}`);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2 flex items-center">
                {ngoProfile.name}
                {ngoProfile.verified && (
                  <svg className="w-8 h-8 text-blue-500 ml-3" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                )}
              </h1>
              <p className="text-gray-600">{ngoProfile.description}</p>
            </div>
            <button
              onClick={handleCreateCampaign}
              className="btn-primary"
            >
              Create New Campaign
            </button>
          </div>
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
                <p className="text-sm text-gray-600">Total Raised</p>
                <p className="text-2xl font-bold text-gray-900">{totalRaised.toFixed(1)} ETH</p>
                <p className="text-sm text-green-600">≈ ${(totalRaised * 2500).toLocaleString()}</p>
              </div>
            </div>
          </div>

          <div className="card">
            <div className="flex items-center">
              <div className="w-12 h-12 bg-secondary-100 rounded-lg flex items-center justify-center mr-4">
                <svg className="w-6 h-6 text-secondary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <div>
                <p className="text-sm text-gray-600">Total Donors</p>
                <p className="text-2xl font-bold text-gray-900">{totalDonors.toLocaleString()}</p>
                <p className="text-sm text-blue-600">Unique contributors</p>
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
                <p className="text-sm text-gray-600">Transparency Score</p>
                <p className="text-2xl font-bold text-gray-900">{ngoProfile.transparencyScore}%</p>
                <p className="text-sm text-yellow-600">Blockchain verified</p>
              </div>
            </div>
          </div>

          <div className="card">
            <div className="flex items-center">
              <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mr-4">
                <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div>
                <p className="text-sm text-gray-600">Active Campaigns</p>
                <p className="text-2xl font-bold text-gray-900">{campaigns.filter(c => c.status === 'Active').length}</p>
                <p className="text-sm text-purple-600">Currently running</p>
              </div>
            </div>
          </div>
        </div>

        {/* Navigation Tabs */}
        <div className="border-b border-gray-200 mb-6">
          <nav className="-mb-px flex space-x-8">
            <button
              onClick={() => setActiveTab('overview')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'overview'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Overview
            </button>
            <button
              onClick={() => setActiveTab('blockchain')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'blockchain'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Blockchain
            </button>
            <button
              onClick={() => setActiveTab('campaigns')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'campaigns'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Campaigns
            </button>
            <button
              onClick={() => setActiveTab('donations')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'donations'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Recent Donations
            </button>
            <button
              onClick={() => setActiveTab('analytics')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'analytics'
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              Analytics
            </button>
          </nav>
        </div>

        {/* Content based on active tab */}
        {activeTab === 'blockchain' && (
          <div className="space-y-8">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Donation Chain */}
              <BlockchainViewer
                ngoId={ngoProfile.id}
                chainType="donations"
              />
              {/* Expenditure Chain */}
              <BlockchainViewer
                ngoId={ngoProfile.id}
                chainType="expenditures"
              />
            </div>

            {/* Latest Transaction Tracker */}
            {latestTransaction && (
              <div className="mt-8">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Latest Transaction</h3>
                <TransactionTracker
                  transactionId={latestTransaction.id}
                  type={latestTransaction.type}
                />
              </div>
            )}
          </div>
        )}

        {activeTab === 'overview' && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Organization Profile */}
              <div className="card">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Organization Profile</h3>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Category:</span>
                    <span className="font-medium">{ngoProfile.category}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Established:</span>
                    <span className="font-medium">{ngoProfile.established}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Location:</span>
                    <span className="font-medium">{ngoProfile.location}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Tax ID:</span>
                    <span className="font-medium font-mono text-sm">{ngoProfile.taxId}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Verification Status:</span>
                    <span className="flex items-center">
                      <span className="text-green-600 font-medium">Verified</span>
                      <svg className="w-4 h-4 text-green-500 ml-1" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                      </svg>
                    </span>
                  </div>
                </div>
              </div>

              {/* Recent Activity */}
              <div className="card">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h3>
                <div className="space-y-3">
                  {recentDonations.slice(0, 3).map((donation) => (
                    <div key={donation.id} className="flex justify-between items-center py-2 border-b border-gray-100 last:border-b-0">
                      <div>
                        <p className="text-sm font-medium text-gray-900">{donation.amount}</p>
                        <p className="text-xs text-gray-500">from {donation.donor}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm text-gray-600">{donation.date}</p>
                        <p className="text-xs text-gray-500">{donation.campaign}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'campaigns' && (
          <div className="space-y-6">
            <div className="flex justify-between items-center">
              <h2 className="text-xl font-semibold text-gray-900">Campaign Management</h2>
              <button onClick={handleCreateCampaign} className="btn-primary">
                Create New Campaign
              </button>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {campaigns.map((campaign) => (
                <div key={campaign.id} className="card">
                  <div className="flex justify-between items-start mb-4">
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900 mb-1">{campaign.title}</h3>
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        campaign.status === 'Active' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-gray-100 text-gray-800'
                      }`}>
                        {campaign.status}
                      </span>
                    </div>
                    <span className="text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded">{campaign.category}</span>
                  </div>
                  
                  <p className="text-gray-600 mb-4">{campaign.description}</p>
                  
                  <div className="space-y-3 mb-4">
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Progress:</span>
                      <span className="font-semibold">{campaign.raised} / {campaign.goal} ({campaign.progress}%)</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-primary-600 h-2 rounded-full" 
                        style={{ width: `${campaign.progress}%` }}
                      ></div>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span className="text-gray-600">Donors: {campaign.donors}</span>
                      <span className="text-gray-600">
                        {campaign.daysLeft > 0 ? `${campaign.daysLeft} days left` : 'Campaign ended'}
                      </span>
                    </div>
                  </div>
                  
                  <div className="flex space-x-3">
                    {campaign.status === 'Active' && campaign.progress === 100 && (
                      <button
                        onClick={() => handleWithdrawFunds(campaign.id)}
                        className="flex-1 btn-secondary"
                      >
                        Withdraw Funds
                      </button>
                    )}
                    <button className="flex-1 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors duration-200">
                      View Details
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'donations' && (
          <div className="space-y-6">
            <h2 className="text-xl font-semibold text-gray-900">Recent Donations</h2>
            
            <div className="card overflow-hidden">
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Donor</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Campaign</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Transaction</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {recentDonations.map((donation) => (
                      <tr key={donation.id} className="hover:bg-gray-50">
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-mono text-gray-900">{donation.donor}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm font-medium text-gray-900">{donation.amount}</div>
                          <div className="text-sm text-gray-500">{donation.usdAmount}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="text-sm text-gray-900">{donation.campaign}</div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {donation.date}
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

        {activeTab === 'analytics' && (
          <div className="space-y-6">
            <h2 className="text-xl font-semibold text-gray-900">Analytics & Reports</h2>
            
            <div className="card">
              <div className="text-center py-12">
                <svg className="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
                <h3 className="text-lg font-medium text-gray-900 mb-2">Advanced Analytics Coming Soon</h3>
                <p className="text-gray-600 mb-4">
                  Detailed charts, donation trends, and performance metrics will be available once integrated with the Go backend.
                </p>
                <div className="space-y-2 text-sm text-gray-500">
                  <p>• Donation trends over time</p>
                  <p>• Geographic distribution of donors</p>
                  <p>• Campaign performance metrics</p>
                  <p>• Blockchain transaction analytics</p>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default NGODashboard;