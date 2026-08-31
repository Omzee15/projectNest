import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Navbar } from '@/components/Navbar';
import {
  Download,
  CheckCircle,
  Grid,
  Lock,
  Terminal,
  ListChecks,
} from 'lucide-react';

// Keep in sync with the files in /public/downloads
const DOWNLOAD_VERSION = '1.0.0';
const DMG_FILE = `TodoBar-macOS-${DOWNLOAD_VERSION}.dmg`;
const DMG_URL = `/downloads/${DMG_FILE}`;
const ZIP_FILE = `TodoBar-macOS-${DOWNLOAD_VERSION}.zip`;
const ZIP_URL = `/downloads/${ZIP_FILE}`;

const features = [
  {
    icon: Lock,
    title: 'Sign in once',
    description:
      'Use your ProjectNest account. Your session is stored in the macOS Keychain and restored on launch.',
  },
  {
    icon: Grid,
    title: 'All your projects & lists',
    description:
      'Pick a project, then a list, straight from the menu bar — no need to open the web app.',
  },
  {
    icon: ListChecks,
    title: 'Add & manage tasks',
    description:
      'Create tasks, tick them complete, or delete them. Everything syncs with your ProjectNest workspace.',
  },
];

const DownloadTodoBar = () => {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />

      <main className="container mx-auto px-4 sm:px-6 py-8 sm:py-16">
        {/* Hero */}
        <div className="text-center max-w-3xl mx-auto mb-12 sm:mb-16">
          <div className="mx-auto w-14 h-14 bg-primary/10 rounded-2xl flex items-center justify-center mb-6">
            <ListChecks className="h-7 w-7 text-primary" />
          </div>
          <h1 className="text-3xl sm:text-4xl md:text-5xl font-bold text-foreground mb-4">
            TodoBar for macOS
          </h1>
          <p className="text-base sm:text-lg md:text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            A native macOS menu bar app for ProjectNest. Capture tasks into your project
            lists in a couple of clicks, without leaving what you're doing.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 justify-center items-center">
            <Button size="lg" asChild className="w-full sm:w-auto">
              <a href={DMG_URL} download>
                <Download className="mr-2 h-4 w-4" />
                Download for macOS (.dmg)
              </a>
            </Button>
            <span className="text-sm text-muted-foreground">
              Universal · macOS 13 Ventura or later · v{DOWNLOAD_VERSION}
            </span>
          </div>
          <p className="text-xs text-muted-foreground mt-3">
            Prefer a zip?{' '}
            <a href={ZIP_URL} download className="underline hover:text-foreground">
              Download the .zip
            </a>{' '}
            instead.
          </p>
        </div>

        {/* Features */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 sm:gap-6 md:gap-8 mb-12 sm:mb-16">
          {features.map((feature, index) => (
            <Card key={index} className="text-center">
              <CardHeader>
                <div className="mx-auto w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                  <feature.icon className="h-6 w-6 text-primary" />
                </div>
                <CardTitle className="text-xl">{feature.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-base text-muted-foreground">{feature.description}</p>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Install instructions */}
        <div className="max-w-3xl mx-auto">
          <h2 className="text-2xl sm:text-3xl font-bold text-foreground mb-6">How to install</h2>

          <ol className="space-y-6">
            <li className="flex gap-4">
              <StepNumber n={1} />
              <div>
                <h3 className="font-semibold text-foreground mb-1">Open the disk image</h3>
                <p className="text-muted-foreground">
                  Click <span className="font-medium text-foreground">Download for macOS</span> above
                  to get{' '}
                  <code className="px-1.5 py-0.5 rounded bg-muted text-sm">{DMG_FILE}</code>, then
                  double-click it in your <span className="font-medium text-foreground">Downloads</span>{' '}
                  folder. A window opens showing the TodoBar icon and an Applications shortcut.
                </p>
              </div>
            </li>

            <li className="flex gap-4">
              <StepNumber n={2} />
              <div>
                <h3 className="font-semibold text-foreground mb-1">
                  Drag TodoBar into Applications
                </h3>
                <p className="text-muted-foreground">
                  In that window, drag the{' '}
                  <code className="px-1.5 py-0.5 rounded bg-muted text-sm">TodoBar</code> icon onto the{' '}
                  <span className="font-medium text-foreground">Applications</span> folder next to it.
                  When it finishes copying, close the window and eject the disk image (click the ⏏
                  next to it in a Finder sidebar), then move{' '}
                  <code className="px-1.5 py-0.5 rounded bg-muted text-sm">{DMG_FILE}</code> to the
                  Trash.
                </p>
              </div>
            </li>

            <li className="flex gap-4">
              <StepNumber n={3} />
              <div>
                <h3 className="font-semibold text-foreground mb-1 flex items-center gap-2">
                  <Terminal className="h-4 w-4" />
                  Remove the quarantine flag
                </h3>
                <p className="text-muted-foreground mb-3">
                  TodoBar is not distributed through the App Store, so macOS quarantines it and shows
                  a message like{' '}
                  <span className="italic">
                    “Apple could not verify ‘TodoBar’ is free of malware.”
                  </span>{' '}
                  Open the <span className="font-medium text-foreground">Terminal</span> app (press ⌘
                  Space, type “Terminal”, Return), paste the command below, press Return, and enter
                  your Mac password when prompted. You only need to do this once.
                </p>
                <pre className="text-sm bg-background rounded-md p-3 overflow-x-auto border border-border">
                  <code>sudo xattr -rd com.apple.quarantine /Applications/TodoBar.app</code>
                </pre>
              </div>
            </li>

            <li className="flex gap-4">
              <StepNumber n={4} />
              <div>
                <h3 className="font-semibold text-foreground mb-1">Open it and sign in</h3>
                <p className="text-muted-foreground">
                  Open <span className="font-medium text-foreground">TodoBar</span> from your
                  Applications folder. A checklist icon appears in your menu bar (top-right of the
                  screen). Click it and sign in with the same email and password you use for
                  ProjectNest — your projects and lists load right away, and your session is
                  remembered.
                </p>
              </div>
            </li>

            <li className="flex gap-4">
              <StepNumber n={5} />
              <div>
                <h3 className="font-semibold text-foreground mb-1">Optional — launch at login</h3>
                <p className="text-muted-foreground">
                  Open{' '}
                  <span className="font-medium text-foreground">
                    System Settings → General → Login Items
                  </span>
                  , click <span className="font-medium text-foreground">+</span>, and add{' '}
                  <code className="px-1.5 py-0.5 rounded bg-muted text-sm mx-1">TodoBar.app</code> so
                  it starts automatically.
                </p>
              </div>
            </li>
          </ol>

          {/* Requirements / notes */}
          <div className="mt-10 rounded-lg bg-muted py-6 px-6">
            <h3 className="font-semibold text-foreground mb-3">Good to know</h3>
            <ul className="space-y-2 text-muted-foreground">
              <li className="flex gap-2">
                <CheckCircle className="h-5 w-5 text-primary shrink-0 mt-0.5" />
                Ships as a standard drag-to-Applications disk image (.dmg). A .zip is also available
                above.
              </li>
              <li className="flex gap-2">
                <CheckCircle className="h-5 w-5 text-primary shrink-0 mt-0.5" />
                Works on both Apple Silicon and Intel Macs (universal build).
              </li>
              <li className="flex gap-2">
                <CheckCircle className="h-5 w-5 text-primary shrink-0 mt-0.5" />
                Requires macOS 13 (Ventura) or later.
              </li>
              <li className="flex gap-2">
                <CheckCircle className="h-5 w-5 text-primary shrink-0 mt-0.5" />
                Runs only in the menu bar — no Dock icon and no window clutter.
              </li>
              <li className="flex gap-2">
                <CheckCircle className="h-5 w-5 text-primary shrink-0 mt-0.5" />
                Talks to the same ProjectNest backend as this website. You need a ProjectNest account.
              </li>
            </ul>
          </div>

          <div className="mt-8 text-center">
            <Button size="lg" asChild>
              <a href={DMG_URL} download>
                <Download className="mr-2 h-4 w-4" />
                Download TodoBar {DOWNLOAD_VERSION}
              </a>
            </Button>
          </div>
        </div>
      </main>
    </div>
  );
};

const StepNumber = ({ n }: { n: number }) => (
  <div className="shrink-0 w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-semibold text-sm">
    {n}
  </div>
);

export default DownloadTodoBar;
