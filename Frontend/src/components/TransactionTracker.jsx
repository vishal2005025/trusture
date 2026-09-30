import React, { useState, useEffect } from 'react';

const TransactionStatus = {
  PENDING: 'pending',
  MINING: 'mining',
  VALIDATING: 'validating',
  COMPLETED: 'completed',
  FAILED: 'failed'
};

const TransactionTracker = ({ transactionId, type }) => {
  const [status, setStatus] = useState(TransactionStatus.PENDING);
  const [details, setDetails] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (!transactionId) return;
    
    const checkStatus = async () => {
      try {
        const response = await fetch(`/api/v1/transactions/${type}/${transactionId}`);
        if (!response.ok) throw new Error('Failed to fetch transaction status');
        
        const data = await response.json();
        setDetails(data);
        setStatus(data.status);

        // Continue polling if not in final state
        return [TransactionStatus.COMPLETED, TransactionStatus.FAILED].includes(data.status);
      } catch (err) {
        setError(err.message);
        return true; // Stop polling on error
      }
    };

    // Initial check
    checkStatus();

    // Poll every 5 seconds until final state
    const interval = setInterval(async () => {
      const shouldStop = await checkStatus();
      if (shouldStop) clearInterval(interval);
    }, 5000);

    return () => clearInterval(interval);
  }, [transactionId, type]);

  const getStatusIcon = () => {
    switch (status) {
      case TransactionStatus.PENDING:
        return (
          <svg className="animate-spin h-5 w-5 text-yellow-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        );
      case TransactionStatus.MINING:
        return (
          <svg className="animate-pulse h-5 w-5 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        );
      case TransactionStatus.VALIDATING:
        return (
          <svg className="h-5 w-5 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
        );
      case TransactionStatus.COMPLETED:
        return (
          <svg className="h-5 w-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
          </svg>
        );
      case TransactionStatus.FAILED:
        return (
          <svg className="h-5 w-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        );
      default:
        return null;
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case TransactionStatus.PENDING:
        return 'text-yellow-700 bg-yellow-50';
      case TransactionStatus.MINING:
        return 'text-blue-700 bg-blue-50';
      case TransactionStatus.VALIDATING:
        return 'text-indigo-700 bg-indigo-50';
      case TransactionStatus.COMPLETED:
        return 'text-green-700 bg-green-50';
      case TransactionStatus.FAILED:
        return 'text-red-700 bg-red-50';
      default:
        return 'text-gray-700 bg-gray-50';
    }
  };

  const getStatusText = () => {
    switch (status) {
      case TransactionStatus.PENDING:
        return 'Transaction Pending';
      case TransactionStatus.MINING:
        return 'Mining Block';
      case TransactionStatus.VALIDATING:
        return 'Validating Transaction';
      case TransactionStatus.COMPLETED:
        return 'Transaction Complete';
      case TransactionStatus.FAILED:
        return 'Transaction Failed';
      default:
        return 'Unknown Status';
    }
  };

  if (error) {
    return (
      <div className="rounded-lg border border-red-200 bg-red-50 p-4">
        <div className="flex items-center text-red-700">
          <svg className="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Error tracking transaction: {error}
        </div>
      </div>
    );
  }

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-4">
      {/* Status Header */}
      <div className={`rounded-lg p-4 ${getStatusColor()}`}>
        <div className="flex items-center">
          {getStatusIcon()}
          <span className="ml-2 font-medium">{getStatusText()}</span>
        </div>
      </div>

      {/* Transaction Details */}
      {details && (
        <div className="mt-4 space-y-4">
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-gray-500">Transaction ID</p>
              <p className="font-mono break-all">{details.transactionId}</p>
            </div>
            <div>
              <p className="text-gray-500">Block Hash</p>
              <p className="font-mono break-all">{details.blockHash || 'Pending...'}</p>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-gray-500">Amount</p>
              <p className="font-medium">{details.amount} ETH</p>
            </div>
            <div>
              <p className="text-gray-500">Timestamp</p>
              <p className="font-medium">
                {details.timestamp ? new Date(details.timestamp).toLocaleString() : 'Pending...'}
              </p>
            </div>
          </div>

          {/* Validation Progress */}
          {details.validators && (
            <div className="mt-4">
              <p className="text-sm text-gray-500 mb-2">Validation Progress</p>
              <div className="flex items-center">
                <div className="flex-1 bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-primary-600 h-2 rounded-full transition-all duration-500"
                    style={{
                      width: `${(details.validators.length / details.requiredValidations) * 100}%`
                    }}
                  ></div>
                </div>
                <span className="ml-2 text-sm text-gray-600">
                  {details.validators.length}/{details.requiredValidations}
                </span>
              </div>
            </div>
          )}

          {/* Blockchain Anchoring Status */}
          {details.polygonAnchor && (
            <div className="mt-4 p-3 bg-purple-50 rounded-lg">
              <div className="flex items-center">
                <svg className="h-5 w-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                <span className="ml-2 text-sm text-purple-700">
                  Anchored to Polygon (Block #{details.polygonAnchor.blockNumber})
                </span>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default TransactionTracker;