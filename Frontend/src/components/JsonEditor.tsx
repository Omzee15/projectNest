import React, { useState, useEffect, useRef } from 'react';
import Editor, { Monaco } from '@monaco-editor/react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  FileJson,
  Check,
  AlertCircle,
  Copy,
  Download,
  Upload,
  RotateCcw,
  Wand2,
} from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

interface JsonEditorProps {
  initialValue?: string;
  onChange?: (value: string) => void;
  readOnly?: boolean;
}

const JsonEditor: React.FC<JsonEditorProps> = ({
  initialValue = '{\n  "example": "Start typing your JSON here"\n}',
  onChange,
  readOnly = false,
}) => {
  const { toast } = useToast();
  const [jsonValue, setJsonValue] = useState(initialValue);
  const [validationError, setValidationError] = useState<string | null>(null);
  const [isValid, setIsValid] = useState(true);
  const editorRef = useRef<any>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    validateJson(jsonValue);
  }, []);

  const validateJson = (value: string) => {
    if (!value.trim()) {
      setValidationError(null);
      setIsValid(true);
      return;
    }

    try {
      JSON.parse(value);
      setValidationError(null);
      setIsValid(true);
    } catch (error: any) {
      const message = error.message || 'Invalid JSON';
      setValidationError(message);
      setIsValid(false);
    }
  };

  const handleEditorChange = (value: string | undefined) => {
    const newValue = value || '';
    setJsonValue(newValue);
    validateJson(newValue);
    onChange?.(newValue);
  };

  const handleEditorDidMount = (editor: any, monaco: Monaco) => {
    editorRef.current = editor;

    // Configure Monaco options
    monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
      validate: true,
      schemas: [],
      allowComments: false,
    });
  };

  const handleFormat = () => {
    if (!editorRef.current) return;

    try {
      const parsed = JSON.parse(jsonValue);
      const formatted = JSON.stringify(parsed, null, 2);
      setJsonValue(formatted);
      validateJson(formatted);
      onChange?.(formatted);
      
      toast({
        title: 'Formatted',
        description: 'JSON has been formatted successfully.',
        duration: 2000,
      });
    } catch (error: any) {
      toast({
        title: 'Cannot format',
        description: 'Please fix JSON errors before formatting.',
        variant: 'destructive',
      });
    }
  };

  const handleMinify = () => {
    try {
      const parsed = JSON.parse(jsonValue);
      const minified = JSON.stringify(parsed);
      setJsonValue(minified);
      validateJson(minified);
      onChange?.(minified);
      
      toast({
        title: 'Minified',
        description: 'JSON has been minified successfully.',
        duration: 2000,
      });
    } catch (error: any) {
      toast({
        title: 'Cannot minify',
        description: 'Please fix JSON errors before minifying.',
        variant: 'destructive',
      });
    }
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(jsonValue);
      toast({
        title: 'Copied',
        description: 'JSON copied to clipboard.',
        duration: 2000,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to copy to clipboard.',
        variant: 'destructive',
      });
    }
  };

  const handleDownload = () => {
    try {
      const blob = new Blob([jsonValue], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `data-${Date.now()}.json`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      
      toast({
        title: 'Downloaded',
        description: 'JSON file has been downloaded.',
        duration: 2000,
      });
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to download JSON.',
        variant: 'destructive',
      });
    }
  };

  const handleUpload = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      try {
        // Validate it's valid JSON
        JSON.parse(content);
        const formatted = JSON.stringify(JSON.parse(content), null, 2);
        setJsonValue(formatted);
        validateJson(formatted);
        onChange?.(formatted);
        
        toast({
          title: 'Loaded',
          description: 'JSON file has been loaded successfully.',
          duration: 2000,
        });
      } catch (error) {
        toast({
          title: 'Error',
          description: 'Invalid JSON file.',
          variant: 'destructive',
        });
      }
    };
    reader.readAsText(file);
    
    // Reset input so the same file can be uploaded again
    event.target.value = '';
  };

  const handleClear = () => {
    const emptyJson = '{\n  \n}';
    setJsonValue(emptyJson);
    validateJson(emptyJson);
    onChange?.(emptyJson);
  };

  return (
    <div className="h-full flex flex-col bg-gray-50">
      <Card className="flex-1 flex flex-col m-4 overflow-hidden">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <FileJson className="w-5 h-5 text-blue-500" />
              <CardTitle>JSON Editor & Validator</CardTitle>
            </div>
            <div className="flex items-center gap-2">
              {isValid ? (
                <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
                  <Check className="w-3 h-3 mr-1" />
                  Valid JSON
                </Badge>
              ) : (
                <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200">
                  <AlertCircle className="w-3 h-3 mr-1" />
                  Invalid JSON
                </Badge>
              )}
            </div>
          </div>
          {validationError && (
            <div className="mt-2 p-2 bg-red-50 border border-red-200 rounded-md flex items-start gap-2">
              <AlertCircle className="w-4 h-4 text-red-500 mt-0.5 flex-shrink-0" />
              <div className="text-sm text-red-700">
                <span className="font-medium">Error: </span>
                {validationError}
              </div>
            </div>
          )}
        </CardHeader>

        <CardContent className="flex-1 flex flex-col p-0 overflow-hidden">
          {/* Toolbar */}
          <div className="px-4 py-2 border-t border-b bg-gray-50 flex flex-wrap gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={handleFormat}
              disabled={!isValid || readOnly}
              title="Format JSON (Beautify)"
            >
              <Wand2 className="w-4 h-4 mr-2" />
              Format
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={handleMinify}
              disabled={!isValid || readOnly}
              title="Minify JSON"
            >
              <Wand2 className="w-4 h-4 mr-2" />
              Minify
            </Button>
            
            <div className="w-px h-6 bg-gray-300" />
            
            <Button
              variant="outline"
              size="sm"
              onClick={handleCopy}
              title="Copy to clipboard"
            >
              <Copy className="w-4 h-4 mr-2" />
              Copy
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={handleDownload}
              title="Download JSON file"
            >
              <Download className="w-4 h-4 mr-2" />
              Download
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={handleUpload}
              disabled={readOnly}
              title="Upload JSON file"
            >
              <Upload className="w-4 h-4 mr-2" />
              Upload
            </Button>
            
            <div className="w-px h-6 bg-gray-300" />
            
            <Button
              variant="outline"
              size="sm"
              onClick={handleClear}
              disabled={readOnly}
              title="Clear editor"
            >
              <RotateCcw className="w-4 h-4 mr-2" />
              Clear
            </Button>
          </div>

          {/* Editor */}
          <div className="flex-1 relative">
            <Editor
              height="100%"
              defaultLanguage="json"
              value={jsonValue}
              onChange={handleEditorChange}
              onMount={handleEditorDidMount}
              theme="vs-light"
              options={{
                minimap: { enabled: false },
                scrollBeyondLastLine: false,
                fontSize: 20,
                lineNumbers: 'on',
                roundedSelection: false,
                readOnly: readOnly,
                automaticLayout: true,
                tabSize: 2,
                wordWrap: 'on',
                formatOnPaste: true,
                formatOnType: true,
                folding: true,
                bracketPairColorization: {
                  enabled: true,
                },
              }}
            />
          </div>

          {/* File input (hidden) */}
          <input
            ref={fileInputRef}
            type="file"
            accept=".json,application/json"
            onChange={handleFileChange}
            className="hidden"
          />
        </CardContent>
      </Card>
    </div>
  );
};

export default JsonEditor;
