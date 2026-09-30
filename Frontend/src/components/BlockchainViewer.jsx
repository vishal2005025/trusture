import React, { useState, useEffect } from 'react';

const BlockchainViewer = ({ ngoId, chainType }) => {
  const [blocks, setBlocks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchBlockchain();
    // Poll for updates every 30 seconds
    const interval = setInterval(fetchBlockchain, 30000);
    return () => clearInterval(interval);
  }, [ngoId, chainType]);

  const fetchBlockchain = async () => {
    try {
      const response = await fetch(`/api/v1/blockchain/${chainType}/${ngoId}/blocks`);
      if (!response.ok) throw new Error('Failed to fetch blockchain data');
      const data = await response.json();
      setBlocks(data.blocks);
      setLoading(false);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  if (loading) return (
    <div className="flex justify-center items-center h-64">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>
  );

  if (error) return (
    <div className="text-red-600 p-4 text-center">
      Error loading blockchain data: {error}
    </div>
  );

  return (
    <div className="bg-white rounded-lg shadow-sm">
      <div className="p-4 border-b border-gray-200">
        <h3 className="text-lg font-semibold text-gray-900">
          {chainType === 'donations' ? 'Donation' : 'Expenditure'} Blockchain
        </h3>
        <p className="text-sm text-gray-600">
          Real-time view of verified transactions
        </p>
      </div>

      <div className="p-4">
        <div className="space-y-4">
          {blocks.map((block, index) => (
            <div
              key={block.hash}
              className="relative pl-8 pb-4 group hover:bg-gray-50 rounded-lg transition-colors duration-150"
            >
              {/* Connection line */}
              {index < blocks.length - 1 && (
                <div className="absolute left-4 top-8 bottom-0 w-0.5 bg-gray-300"></div>
              )}

              {/* Block dot */}
              <div className="absolute left-2 top-2 w-4 h-4 rounded-full bg-primary-100 border-2 border-primary-500 group-hover:bg-primary-500 transition-colors duration-150"></div>

              {/* Block content */}
              <div className="bg-white border border-gray-200 rounded-lg p-4 group-hover:border-primary-500 transition-colors duration-150">
                <div className="flex items-start justify-between mb-2">
                  <div>
                    <span className="text-xs font-medium text-gray-500">Block #{block.index}</span>
                    <h4 className="text-sm font-semibold text-gray-900">{block.type}</h4>
                  </div>
                  <div className="flex items-center space-x-2">
                    {block.validated && (
                      <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                        Validated
                      </span>
                    )}
                    {block.anchored && (
                      <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800">
                        Polygon
                      </span>
                    )}
                  </div>
                </div>

                <div className="space-y-1 text-sm">
                  <p className="font-mono text-gray-600 text-xs break-all">
                    Hash: {block.hash}
                  </p>
                  <p className="font-mono text-gray-500 text-xs break-all">
                    Prev: {block.previousHash}
                  </p>
                </div>

                <div className="mt-2 pt-2 border-t border-gray-100">
                  <div className="text-sm text-gray-900">
                    <div className="flex justify-between items-center">
                      <span>Timestamp</span>
                      <span className="font-medium">{new Date(block.timestamp).toLocaleString()}</span>
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span>Validators</span>
                      <span className="font-medium">{block.validators.length}</span>
                    </div>
                    {block.merkleRoot && (
                      <div className="flex justify-between items-center mt-1">
                        <span>Merkle Root</span>
                        <span className="font-mono text-xs">{block.merkleRoot.slice(0, 10)}...</span>
                      </div>
                    )}
                  </div>
                </div>

                <div className="mt-2 pt-2 border-t border-gray-100">
                  <details className="text-sm">
                    <summary className="cursor-pointer text-primary-600 hover:text-primary-700">
                      View Transaction Data
                    </summary>
                    <div className="mt-2 bg-gray-50 rounded p-2 overflow-auto max-h-40">
                      <pre className="text-xs text-gray-700">
                        {JSON.stringify(block.data, null, 2)}
                      </pre>
                    </div>
                  </details>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {blocks.length === 0 && (
        <div className="p-8 text-center text-gray-500">
          No blocks found in this chain yet
        </div>
      )}
    </div>
  );
};

export default BlockchainViewer;