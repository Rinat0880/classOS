import { useQuery } from '@tanstack/react-query';
import { X, Activity, Globe, AppWindow, LogIn, LogOut } from 'lucide-react';
import { logsService } from '../../services/logs';
import type { UserLog } from '../../types';

interface LogsModalProps {
  isOpen: boolean;
  onClose: () => void;
  username?: string;
  deviceName?: string;
}

const LogsModal = ({ isOpen, onClose, username, deviceName }: LogsModalProps) => {
  const { data: logs = [], isLoading } = useQuery({
    queryKey: ['logs', username, deviceName],
    queryFn: () => {
      if (username) {
        return logsService.getByUsername(username, 100);
      }
      if (deviceName) {
        return logsService.getByDevice(deviceName, 100);
      }
      return Promise.resolve([]);
    },
    enabled: isOpen && (!!username || !!deviceName),
  });

  const getLogIcon = (log: UserLog) => {
    if (log.log_type === 'system') {
      return log.action.includes('Start') ? (
        <LogIn className="w-4 h-4 text-green-600" />
      ) : (
        <LogOut className="w-4 h-4 text-red-600" />
      );
    }
    if (log.log_type === 'process') {
      return <AppWindow className="w-4 h-4 text-purple-600" />;
    }
    return <Globe className="w-4 h-4 text-orange-600" />;
  };

  const getLogTypeColor = (type: string) => {
    switch (type) {
      case 'system':
        return 'bg-blue-100 text-blue-800';
      case 'process':
        return 'bg-purple-100 text-purple-800';
      case 'browser':
        return 'bg-orange-100 text-orange-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    return new Intl.DateTimeFormat('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    }).format(date);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-6xl max-h-[90vh] flex flex-col">
        <div className="flex items-center justify-between p-6 border-b">
          <div>
            <h2 className="text-xl font-bold text-gray-900 flex items-center gap-2">
              <Activity className="w-6 h-6 text-blue-600" />
              Activity Logs
            </h2>
            <p className="text-sm text-gray-600 mt-1">
              {username && `User: ${username}`}
              {deviceName && `Device: ${deviceName}`}
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <div className="flex-1 overflow-auto p-6">
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            </div>
          ) : logs.length === 0 ? (
            <div className="text-center py-12">
              <Activity className="w-12 h-12 text-gray-400 mx-auto mb-3" />
              <p className="text-gray-500">No activity logs found</p>
            </div>
          ) : (
            <div className="space-y-3">
              {logs.map((log) => (
                <div
                  key={log.id}
                  className="bg-gray-50 rounded-lg p-4 border border-gray-200 hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-start gap-3">
                    <div className="mt-1">{getLogIcon(log)}</div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-2">
                        <span
                          className={`px-2 py-1 rounded text-xs font-medium ${getLogTypeColor(
                            log.log_type
                          )}`}
                        >
                          {log.log_type}
                        </span>
                        <span className="text-xs text-gray-500">
                          {formatTimestamp(log.timestamp)}
                        </span>
                      </div>
                      <div className="flex items-start gap-2">
                        <div className="flex-1">
                          {log.program && (
                            <p className="text-sm font-medium text-gray-900 mb-1">
                              {log.program}
                            </p>
                          )}
                          <p className="text-sm text-gray-600 break-all">{log.action}</p>
                        </div>
                      </div>
                      <div className="flex items-center gap-3 mt-2 text-xs text-gray-500">
                        <span>👤 {log.username}</span>
                        <span>💻 {log.device_name}</span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="p-4 border-t bg-gray-50">
          <div className="flex justify-between items-center">
            <p className="text-sm text-gray-600">
              Showing {logs.length} recent activities
            </p>
            <button
              onClick={onClose}
              className="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LogsModal;
