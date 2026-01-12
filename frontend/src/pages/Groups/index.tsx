import { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Plus } from 'lucide-react';
import toast from 'react-hot-toast';
import { groupsService } from '../../services/groups';
import { usersService } from '../../services/users';
import GroupsTable from '../../components/common/groupsTable';
import GroupModal from '../../components/common/GroupsModal';
import DeleteConfirmDialog from '../../components/common/DeleteConfirmDialog';
import GroupDetailsModal from '../../components/common/GroupDetailsModal';
import type { Group, GroupWithUserCount, CreateGroupInput, UpdateGroupInput } from '../../types';

const Groups = () => {
  const queryClient = useQueryClient();
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isDetailsModalOpen, setIsDetailsModalOpen] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<GroupWithUserCount | null>(null);

  const { data: groups = [], isLoading: loadingGroups, error: groupsError } = useQuery({
    queryKey: ['groups'],
    queryFn: groupsService.getAll,
  });

  const { data: allUsers = [], error: usersError } = useQuery({
    queryKey: ['users'],
    queryFn: usersService.getAll,
  });

  const { data: groupUsers = [], isLoading: loadingGroupUsers } = useQuery({
    queryKey: ['group-users', selectedGroup?.id],
    queryFn: () => groupsService.getUsers(selectedGroup!.id),
    enabled: !!selectedGroup && isDetailsModalOpen,
  });

  const groupsWithUserCount: GroupWithUserCount[] = useMemo(() => {
    if (!Array.isArray(groups) || !Array.isArray(allUsers)) {
      return [];
    }
    return groups.map((group) => ({
      ...group,
      userCount: allUsers.filter((user) => user.group_id === group.id).length,
    }));
  }, [groups, allUsers]);

  const createMutation = useMutation({
    mutationFn: (data: CreateGroupInput) => groupsService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      toast.success('Group created successfully!');
      setIsCreateModalOpen(false);
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to create group');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateGroupInput }) =>
      groupsService.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      toast.success('Group updated successfully!');
      setIsEditModalOpen(false);
      setSelectedGroup(null);
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to update group');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => groupsService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      queryClient.invalidateQueries({ queryKey: ['users'] });
      toast.success('Group deleted successfully!');
      setIsDeleteDialogOpen(false);
      setSelectedGroup(null);
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to delete group');
    },
  });

  const handleCreate = (data: CreateGroupInput) => {
    createMutation.mutate(data);
  };

  const handleEdit = (group: GroupWithUserCount) => {
    setSelectedGroup(group);
    setIsEditModalOpen(true);
  };

  const handleUpdate = (data: UpdateGroupInput) => {
    if (selectedGroup) {
      updateMutation.mutate({ id: selectedGroup.id, data });
    }
  };

  const handleDelete = (group: GroupWithUserCount) => {
    setSelectedGroup(group);
    setIsDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (selectedGroup) {
      deleteMutation.mutate(selectedGroup.id);
    }
  };

  const handleViewDetails = (group: GroupWithUserCount) => {
    setSelectedGroup(group);
    setIsDetailsModalOpen(true);
  };

  return (
    <div className="space-y-6">
      {(groupsError || usersError) && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800 font-medium">Error loading data</p>
          <p className="text-sm text-red-600 mt-1">
            {(groupsError as any)?.message || (usersError as any)?.message || 'Please try again'}
          </p>
        </div>
      )}

      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Groups</h1>
          <p className="text-gray-600 mt-1">Manage groups and their members</p>
        </div>
        <button
          onClick={() => setIsCreateModalOpen(true)}
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          <Plus className="w-5 h-5" />
          Add Group
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm font-medium text-gray-600">Total Groups</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">
            {loadingGroups ? '...' : groups.length}
          </p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm font-medium text-gray-600">Total Users</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">{allUsers.length}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm font-medium text-gray-600">Users without Group</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">
            {allUsers.filter((u) => u.group_id === 0).length}
          </p>
        </div>
      </div>

      <GroupsTable
        groups={groupsWithUserCount}
        loading={loadingGroups}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onViewDetails={handleViewDetails}
      />

      <GroupModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={handleCreate}
        isSubmitting={createMutation.isPending}
      />

      <GroupModal
        isOpen={isEditModalOpen}
        onClose={() => {
          setIsEditModalOpen(false);
          setSelectedGroup(null);
        }}
        onSubmit={handleUpdate}
        group={selectedGroup}
        isSubmitting={updateMutation.isPending}
      />

      <DeleteConfirmDialog
        isOpen={isDeleteDialogOpen}
        onClose={() => {
          setIsDeleteDialogOpen(false);
          setSelectedGroup(null);
        }}
        onConfirm={confirmDelete}
        title="Delete Group"
        message={`Are you sure you want to delete "${selectedGroup?.name}"? Users in this group will be moved to "No Group".`}
        isDeleting={deleteMutation.isPending}
      />

      <GroupDetailsModal
        isOpen={isDetailsModalOpen}
        onClose={() => {
          setIsDetailsModalOpen(false);
          setSelectedGroup(null);
        }}
        group={selectedGroup}
        users={groupUsers}
        loadingUsers={loadingGroupUsers}
      />
    </div>
  );
};

export default Groups;