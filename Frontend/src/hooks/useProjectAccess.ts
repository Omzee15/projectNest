import { useEffect, useState } from 'react';
import { apiService } from '@/services/api';
import { ProjectMember } from '@/types';
import { getCurrentUserUid } from '@/utils/auth';

interface ProjectAccess {
  canWrite: boolean;
  isOwner: boolean;
  role: string;
  isMember: boolean;
  loading: boolean;
}

/**
 * Hook to check the current user's access level in a project
 * Returns an object with the user's access permissions
 */
export function useProjectAccess(projectId: string): ProjectAccess {
  const [access, setAccess] = useState<ProjectAccess>({
    canWrite: false,
    isOwner: false,
    role: 'guest',
    isMember: false,
    loading: true,
  });

  useEffect(() => {
    const checkAccess = async () => {
      try {
        const currentUserUid = getCurrentUserUid();
        if (!currentUserUid) {
          setAccess({
            canWrite: false,
            isOwner: false,
            role: 'guest',
            isMember: false,
            loading: false,
          });
          return;
        }

        const response = await apiService.getProjectMembers(projectId);
        const members = response.data || [];

        const currentMember = members.find(
          (member: ProjectMember) => member.user_uid === currentUserUid
        );

        if (!currentMember) {
          setAccess({
            canWrite: false,
            isOwner: false,
            role: 'guest',
            isMember: false,
            loading: false,
          });
          return;
        }

        setAccess({
          canWrite: currentMember.can_write || currentMember.role === 'owner',
          isOwner: currentMember.role === 'owner',
          role: currentMember.role,
          isMember: true,
          loading: false,
        });
      } catch (error) {
        console.error('Failed to check project access:', error);
        setAccess((prev) => ({ ...prev, loading: false }));
      }
    };

    checkAccess();
  }, [projectId]);

  return access;
}
