import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Plus } from 'lucide-react';
import toast from 'react-hot-toast';
import { usersService } from '../../services/users';
import { groupsService } from '../../services/groups';
import UserTable from '../../components/common/UserTable';
import UserModal from '../../components/common/UserModal';
import ChangePasswordModal from '../../components/common/ChangePasswordModal';
import DeleteConfirmDialog from '../../components/common/DeleteConfirmDialog';
import type { User, CreateUserInput, UpdateUserInput } from '../../types';

const Users = () => {
  const queryClient = useQueryClient();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const { data: users = [], isLoading: usersLoading } = useQuery({
    queryKey: ['users'],
    queryFn: usersService.getAll,
  });

  const { data: groups = [] } = useQuery({
    queryKey: ['groups'],
    queryFn: groupsService.getAll,
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateUserInput) => {
      const group = groups.find(g => g.id === data.group_id);
      return groupsService.createUser(data.group_id, {
        name: data.name,
        username: data.username,
        password: data.password,
        role: data.role,
        group_name: group?.name || '',
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      setIsCreateModalOpen(false);
      toast.success('User created successfully');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to create user');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserInput }) =>
      usersService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      setIsEditModalOpen(false);
      setSelectedUser(null);
      toast.success('User updated successfully');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to update user');
    },
  });

  const changePasswordMutation = useMutation({
    mutationFn: ({ id, password }: { id: number; password: string }) =>
      usersService.changePassword(id, password),
    onSuccess: () => {
      setIsPasswordModalOpen(false);
      setSelectedUser(null);
      toast.success('Password changed successfully');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to change password');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => usersService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      setIsDeleteDialogOpen(false);
      setSelectedUser(null);
      toast.success('User deleted successfully');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to delete user');
    },
  });

  const handleCreate = (data: any) => {
    createMutation.mutate(data as CreateUserInput);
  };

  const handleEdit = (data: any) => {
    if (!selectedUser) return;
    
    const updateData: UpdateUserInput = {};
    if (data.name && data.name !== selectedUser.name) updateData.name = data.name;
    if (data.username && data.username !== selectedUser.username) updateData.username = data.username;
    if (data.role && data.role !== selectedUser.role) updateData.role = data.role;
    
    if (data.group_id && data.group_id !== selectedUser.group_id) {
      updateData.group_id = data.group_id;
      const group = groups.find(g => g.id === data.group_id);
      if (group) {
        updateData.group_name = group.name;
      }
    }

    if (Object.keys(updateData).length === 0) {
      toast.error('No changes detected');
      return;
    }

    updateMutation.mutate({ id: selectedUser.id, data: updateData });
  };

  const handleChangePassword = (password: string) => {
    if (!selectedUser) return;
    changePasswordMutation.mutate({ id: selectedUser.id, password });
  };

  const handleDelete = () => {
    if (!selectedUser) return;
    deleteMutation.mutate(selectedUser.id);
  };

  const handleEditClick = (user: User) => {
    setSelectedUser(user);
    setIsEditModalOpen(true);
  };

  const handlePasswordClick = (user: User) => {
    setSelectedUser(user);
    setIsPasswordModalOpen(true);
  };

  const handleDeleteClick = (user: User) => {
    setSelectedUser(user);
    setIsDeleteDialogOpen(true);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Users</h1>
          <p className="text-gray-600 mt-1">Manage users and their permissions</p>
        </div>
        <button
          onClick={() => setIsCreateModalOpen(true)}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <Plus className="w-5 h-5" />
          Add User
        </button>
      </div>

      <UserTable
        users={users}
        loading={usersLoading}
        showActions
        onEdit={handleEditClick}
        onDelete={handleDeleteClick}
        onChangePassword={handlePasswordClick}
      />

      <UserModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={handleCreate}
        groups={groups}
        isLoading={createMutation.isPending}
      />

      <UserModal
        isOpen={isEditModalOpen}
        onClose={() => {
          setIsEditModalOpen(false);
          setSelectedUser(null);
        }}
        onSubmit={handleEdit}
        user={selectedUser}
        groups={groups}
        isLoading={updateMutation.isPending}
      />

      <ChangePasswordModal
        isOpen={isPasswordModalOpen}
        onClose={() => {
          setIsPasswordModalOpen(false);
          setSelectedUser(null);
        }}
        onSubmit={handleChangePassword}
        user={selectedUser}
        isLoading={changePasswordMutation.isPending}
      />

      <DeleteConfirmDialog
        isOpen={isDeleteDialogOpen}
        onClose={() => {
          setIsDeleteDialogOpen(false);
          setSelectedUser(null);
        }}
        onConfirm={handleDelete}
        title="Delete User"
        message={
          selectedUser
            ? `Are you sure you want to delete user "${selectedUser.name}" (@${selectedUser.username})? This action cannot be undone.`
            : ''
        }
        isLoading={deleteMutation.isPending}
      />
    </div>
  );
};

export default Users;
