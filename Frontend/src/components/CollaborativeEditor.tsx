import React, { useEffect, useRef } from 'react';
import { useEditor, EditorContent, Editor } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import {
  Bold,
  Italic,
  Strikethrough,
  Code,
  Heading1,
  Heading2,
  Heading3,
  List,
  ListOrdered,
  Undo,
  Redo,
  Quote,
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface CollaborativeEditorProps {
  documentId: string;
  userName: string;
  userColor?: string;
  onSave?: (content: string, html: string) => void;
  autoSave?: boolean;
  autoSaveInterval?: number;
  placeholder?: string;
  readOnly?: boolean;
  className?: string;
}

export const CollaborativeEditor: React.FC<CollaborativeEditorProps> = ({
  documentId,
  userName,
  userColor = '#3B82F6',
  onSave,
  autoSave = true,
  autoSaveInterval = 3000,
  placeholder = 'Start typing...',
  readOnly = false,
  className,
}) => {
  const saveTimeoutRef = useRef<NodeJS.Timeout>();

  const editor = useEditor({
    editable: !readOnly,
    extensions: [
      StarterKit.configure({
        history: true,
      }),
    ],
    content: '<p>Start typing your collaborative note here...</p>',
    onUpdate: ({ editor }) => {
      if (autoSave && onSave) {
        // Debounce auto-save
        if (saveTimeoutRef.current) {
          clearTimeout(saveTimeoutRef.current);
        }
        saveTimeoutRef.current = setTimeout(() => {
          const json = editor.getJSON();
          const html = editor.getHTML();
          onSave(JSON.stringify(json), html);
        }, autoSaveInterval);
      }
    },
  });

  useEffect(() => {
    return () => {
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current);
      }
    };
  }, []);

  if (!editor) {
    return <div className="animate-pulse h-64 bg-gray-100 rounded-md" />;
  }

  return (
    <div className={cn('border rounded-md overflow-hidden', className)}>
      {!readOnly && <MenuBar editor={editor} />}
      <EditorContent
        editor={editor}
        className={cn(
          'prose prose-sm max-w-none p-4 min-h-[300px] focus:outline-none',
          readOnly && 'bg-gray-50'
        )}
      />
    </div>
  );
};

interface MenuBarProps {
  editor: Editor;
}

const MenuBar: React.FC<MenuBarProps> = ({ editor }) => {
  return (
    <div className="border-b bg-gray-50 p-2 flex flex-wrap gap-1">
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleBold().run()}
        className={cn(editor.isActive('bold') && 'bg-gray-200')}
        type="button"
      >
        <Bold className="w-4 h-4" />
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleItalic().run()}
        className={cn(editor.isActive('italic') && 'bg-gray-200')}
        type="button"
      >
        <Italic className="w-4 h-4" />
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleStrike().run()}
        className={cn(editor.isActive('strike') && 'bg-gray-200')}
        type="button"
      >
        <Strikethrough className="w-4 h-4" />
      </Button>
      
      <Separator orientation="vertical" className="mx-1 h-8" />
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
        className={cn(editor.isActive('heading', { level: 1 }) && 'bg-gray-200')}
        type="button"
      >
        <Heading1 className="w-4 h-4" />
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
        className={cn(editor.isActive('heading', { level: 2 }) && 'bg-gray-200')}
        type="button"
      >
        <Heading2 className="w-4 h-4" />
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
        className={cn(editor.isActive('heading', { level: 3 }) && 'bg-gray-200')}
        type="button"
      >
        <Heading3 className="w-4 h-4" />
      </Button>
      
      <Separator orientation="vertical" className="mx-1 h-8" />
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleBulletList().run()}
        className={cn(editor.isActive('bulletList') && 'bg-gray-200')}
        type="button"
      >
        <List className="w-4 h-4" />
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleOrderedList().run()}
        className={cn(editor.isActive('orderedList') && 'bg-gray-200')}
        type="button"
      >
        <ListOrdered className="w-4 h-4" />
      </Button>
      
      <Separator orientation="vertical" className="mx-1 h-8" />
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleCode().run()}
        className={cn(editor.isActive('code') && 'bg-gray-200')}
        type="button"
      >
        <Code className="w-4 h-4" />
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().toggleBlockquote().run()}
        className={cn(editor.isActive('blockquote') && 'bg-gray-200')}
        type="button"
      >
        <Quote className="w-4 h-4" />
      </Button>
      
      <Separator orientation="vertical" className="mx-1 h-8" />
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().undo().run()}
        disabled={!editor.can().undo()}
        type="button"
      >
        <Undo className="w-4 h-4" />
      </Button>
      
      <Button
        variant="ghost"
        size="sm"
        onClick={() => editor.chain().focus().redo().run()}
        disabled={!editor.can().redo()}
        type="button"
      >
        <Redo className="w-4 h-4" />
      </Button>
    </div>
  );
};
