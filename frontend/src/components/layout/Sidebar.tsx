import { NavLink } from 'react-router-dom';
import { LayoutDashboard, Users, UsersRound, Settings, LogOut, Monitor } from 'lucide-react';
import { ROUTES } from '../../constants';
import { useAuthStore } from '../../store/authStore';

const Sidebar = () => {
  const logout = useAuthStore((state) => state.logout);

  const navItems = [
    {
      name: 'Dashboard',
      path: ROUTES.DASHBOARD,
      icon: LayoutDashboard,
    },
    {
      name: 'Devices',
      path: ROUTES.DEVICES,
      icon: Monitor,
    },
    {
      name: 'Groups',
      path: ROUTES.GROUPS,
      icon: UsersRound,
    },
    {
      name: 'Users',
      path: ROUTES.USERS,
      icon: Users,
    },
    {
      name: 'Settings',
      path: ROUTES.SETTINGS,
      icon: Settings,
    },
  ];

  return (
    <div className="flex flex-col w-64 bg-gray-800">
      <div className="flex items-center justify-center h-16 bg-gray-900">
        <span className="text-white text-lg font-semibold">ClassOS</span>
      </div>

      <nav className="flex-1 px-2 py-4 space-y-1">
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) =>
              `flex items-center px-4 py-3 text-gray-300 rounded-lg transition-colors ${
                isActive
                  ? 'bg-gray-900 text-white'
                  : 'hover:bg-gray-700 hover:text-white'
              }`
            }
          >
            <item.icon className="w-5 h-5 mr-3" />
            <span className="text-sm font-medium">{item.name}</span>
          </NavLink>
        ))}
      </nav>

      <div className="px-2 py-4 border-t border-gray-700">
        <button
          onClick={logout}
          className="flex items-center w-full px-4 py-3 text-gray-300 rounded-lg hover:bg-gray-700 hover:text-white transition-colors"
        >
          <LogOut className="w-5 h-5 mr-3" />
          <span className="text-sm font-medium">Logout</span>
        </button>
      </div>
    </div>
  );
};

export default Sidebar;
