import { useQuery } from '@tanstack/react-query';
import { Users, UserCheck, Shield, FolderOpen } from 'lucide-react';
import { usersService } from '../../services/users';
import UserTable from '../../components/common/UserTable';
import type { DashboardStats } from '../../types';
import { useMemo } from 'react';

const Dashboard = () => {
  // Получаем пользователей
  const { data: users = [], isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: usersService.getAll,
  });

  // Вычисляем статистику из данных
  const stats: DashboardStats = useMemo(() => {
    const totalUsers = users.length;
    const adminUsers = users.filter((u) => u.role === 'admin').length;
    const clientUsers = users.filter((u) => u.role === 'client').length;
    const activeGroups = new Set(users.filter((u) => u.group_name).map((u) => u.group_name))
      .size;
    const onlineUsers = 0; // Заглушка

    return {
      totalUsers,
      adminUsers,
      clientUsers,
      activeGroups,
      onlineUsers,
    };
  }, [users]);

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600 mt-1">Overview of all users and their status</p>
      </div>

      {/* Карточки статистики */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Total Users */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Total Users</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {isLoading ? '...' : stats.totalUsers}
              </p>
            </div>
            <div className="bg-blue-100 p-3 rounded-full">
              <Users className="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </div>

        {/* Admin Users */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Admin Users</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {isLoading ? '...' : stats.adminUsers}
              </p>
            </div>
            <div className="bg-purple-100 p-3 rounded-full">
              <Shield className="w-6 h-6 text-purple-600" />
            </div>
          </div>
        </div>

        {/* Client Users */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Client Users</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {isLoading ? '...' : stats.clientUsers}
              </p>
            </div>
            <div className="bg-green-100 p-3 rounded-full">
              <UserCheck className="w-6 h-6 text-green-600" />
            </div>
          </div>
        </div>

        {/* Active Groups */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Active Groups</p>
              <p className="text-3xl font-bold text-gray-900 mt-2">
                {isLoading ? '...' : stats.activeGroups}
              </p>
            </div>
            <div className="bg-orange-100 p-3 rounded-full">
              <FolderOpen className="w-6 h-6 text-orange-600" />
            </div>
          </div>
        </div>
      </div>

      {/* Таблица пользователей */}
      <UserTable users={users} loading={isLoading} />
    </div>
  );
};

export default Dashboard;