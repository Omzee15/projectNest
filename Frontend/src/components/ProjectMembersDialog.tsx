import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Users } from 'lucide-react';
import { ProjectMembersList } from './ProjectMembersList';

interface ProjectMembersDialogProps {
  projectId: string;
  trigger?: React.ReactNode;
}

export function ProjectMembersDialog({
  projectId,
  trigger,
}: ProjectMembersDialogProps) {
  const [open, setOpen] = useState(false);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button variant="outline" size="sm" className="flex items-center gap-2">
            <Users className="h-4 w-4" />
            Members
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[700px] md:max-w-[750px]">
        <DialogHeader>
          <DialogTitle>Project Members</DialogTitle>
          <DialogDescription>
            Manage team members and their roles for this project.
          </DialogDescription>
        </DialogHeader>
        <div className="py-4">
          <ProjectMembersList projectId={projectId} />
        </div>
      </DialogContent>
    </Dialog>
  );
}
