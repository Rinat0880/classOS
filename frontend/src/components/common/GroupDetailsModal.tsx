import { X, Users, FileText } from 'lucide-react';
import type { Group, User } from '../../types';

interface GroupDetailsModalProps {
  isOpen: boolean;
  onClose: () => void;
  group: Group | null;
  users: User[];
  loadingUsers?: boolean;
}

const GroupDetailsModal = ({
  isOpen,
  onClose,
  group,
  users,
  loadingUsers = false,
}: GroupDetailsModalProps) => {
  if (!isOpen || !group) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black bg-opacity-50" onClick={onClose}></div>

      <div className="relative bg-white rounded-lg shadow-xl w-full max-w-3xl mx-4 max-h-[90vh] overflow-y-auto z-10">
        <div className="sticky top-0 bg-white flex items-center justify-between p-6 border-b border-gray-200">
          <h3 className="text-lg font-semibold text-gray-900">Group Details: {group.name}</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6 space-y-6">
          <div className="bg-gray-50 p-4 rounded-lg">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-gray-600">Group ID</p>
                <p className="text-lg font-semibold text-gray-900">{group.id}</p>
              </div>
              <div>
                <p className="text-sm text-gray-600">Total Users</p>
                <p className="text-lg font-semibold text-gray-900">{users.length}</p>
              </div>
            </div>
          </div>

          <div>
            <div className="flex items-center gap-2 mb-4">
              <Users className="w-5 h-5 text-gray-600" />
              <h4 className="text-md font-semibold text-gray-900">Users in this group</h4>
            </div>
            
            {loadingUsers ? (
              <div className="text-center py-4 text-gray-500">Loading users...</div>
            ) : users.length === 0 ? (
              <div className="text-center py-8 text-gray-500 bg-gray-50 rounded-lg">
                No users in this group
              </div>
            ) : (
              <div className="border border-gray-200 rounded-lg overflow-hidden">
                <table className="w-full">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Name
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Username
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                        Role
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {users.map((user) => (
                      <tr key={user.id} className="hover:bg-gray-50">
                        <td className="px-4 py-3 text-sm text-gray-900">{user.name}</td>
                        <td className="px-4 py-3 text-sm text-gray-600">{user.username}</td>
                        <td className="px-4 py-3">
                          <span
                            className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                              user.role === 'admin'
                                ? 'bg-purple-100 text-purple-800'
                                : 'bg-blue-100 text-blue-800'
                            }`}
                          >
                            {user.role}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          <div>
            <div className="flex items-center gap-2 mb-4">
              <FileText className="w-5 h-5 text-gray-600" />
              <h4 className="text-md font-semibold text-gray-900">Whitelist</h4>
            </div>
            
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-6 text-center">
              <p className="text-yellow-800 font-medium">🚧 Coming Soon</p>
              <p className="text-sm text-yellow-700 mt-2">
                Whitelist management will be implemented in Phase 6.5
              </p>
            </div>
          </div>
        </div>

        <div className="sticky bottom-0 bg-gray-50 px-6 py-4 border-t border-gray-200">
          <button
            onClick={onClose}
            className="w-full px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default GroupDetailsModal;