import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Monitor, Activity, Search } from 'lucide-react';
import { devicesService } from '../../services/devices';
import LogsModal from '../../components/common/LogsModal';
import type { DeviceStatus } from '../../types';

const Devices = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedDevice, setSelectedDevice] = useState<string | null>(null);

  const { data: devices = [], isLoading, refetch } = useQuery({
    queryKey: ['devices', 'online'],
    queryFn: devicesService.getOnline,
    refetchInterval: 30000,
  });

  const filteredDevices = devices.filter(
    (device) =>
      device.device_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      device.username.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const isOnline = (lastHeartbeat: string) => {
    const lastUpdate = new Date(lastHeartbeat);
    const now = new Date();
    const diffMinutes = (now.getTime() - lastUpdate.getTime()) / 1000 / 60;
    return diffMinutes < 2;
  };

  const formatLastSeen = (timestamp: string) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (diffSeconds < 60) return 'Just now';
    if (diffSeconds < 3600) return `${Math.floor(diffSeconds / 60)}m ago`;
    if (diffSeconds < 86400) return `${Math.floor(diffSeconds / 3600)}h ago`;
    return date.toLocaleDateString();
  };

  useEffect(() => {
    const interval = setInterval(() => {
      refetch();
    }, 30000);
    return () => clearInterval(interval);
  }, [refetch]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Online Devices</h1>
        <p className="text-gray-600 mt-1">Monitor active devices in real-time</p>
      </div>

      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-3">
            <div className="bg-green-100 p-3 rounded-full">
              <Monitor className="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm font-medium text-gray-600">Online Devices</p>
              <p className="text-3xl font-bold text-gray-900">
                {isLoading ? '...' : filteredDevices.filter((d) => isOnline(d.last_heartbeat)).length}
              </p>
            </div>
          </div>

          <div className="relative w-64">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search devices or users..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        ) : filteredDevices.length === 0 ? (
          <div className="text-center py-12">
            <Monitor className="w-12 h-12 text-gray-400 mx-auto mb-3" />
            <p className="text-gray-500">No devices found</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-3 px-4 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="text-left py-3 px-4 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Device Name
                  </th>
                  <th className="text-left py-3 px-4 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Username
                  </th>
                  <th className="text-left py-3 px-4 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Last Seen
                  </th>
                  <th className="text-left py-3 px-4 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {filteredDevices.map((device, index) => {
                  const online = isOnline(device.last_heartbeat);
                  return (
                    <tr
                      key={`${device.device_name}-${index}`}
                      className="hover:bg-gray-50 transition-colors"
                    >
                      <td className="py-3 px-4">
                        <div className="flex items-center gap-2">
                          <div
                            className={`w-2 h-2 rounded-full ${
                              online ? 'bg-green-500' : 'bg-gray-400'
                            }`}
                          />
                          <span
                            className={`text-xs font-medium ${
                              online ? 'text-green-600' : 'text-gray-500'
                            }`}
                          >
                            {online ? 'Online' : 'Offline'}
                          </span>
                        </div>
                      </td>
                      <td className="py-3 px-4">
                        <span className="text-sm font-medium text-gray-900">
                          {device.device_name}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <span className="text-sm text-gray-600">{device.username}</span>
                      </td>
                      <td className="py-3 px-4">
                        <span className="text-sm text-gray-500">
                          {formatLastSeen(device.last_heartbeat)}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        <button
                          onClick={() => setSelectedDevice(device.device_name)}
                          className="flex items-center gap-1 px-3 py-1 text-sm text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                        >
                          <Activity className="w-4 h-4" />
                          View Logs
                        </button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <LogsModal
        isOpen={!!selectedDevice}
        onClose={() => setSelectedDevice(null)}
        deviceName={selectedDevice || undefined}
      />
    </div>
  );
};

export default Devices;
