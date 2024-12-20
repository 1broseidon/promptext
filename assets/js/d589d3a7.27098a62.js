"use strict";(self.webpackChunkpromptext_docs=self.webpackChunkpromptext_docs||[]).push([[924],{6475:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>a,contentTitle:()=>o,default:()=>h,frontMatter:()=>r,metadata:()=>i,toc:()=>c});const i=JSON.parse('{"id":"getting-started","title":"Getting Started","description":"Installation","source":"@site/docs/getting-started.md","sourceDirName":".","slug":"/getting-started","permalink":"/promptext/getting-started","draft":false,"unlisted":false,"editUrl":"https://github.com/1broseidon/promptext/tree/main/docs/docs/getting-started.md","tags":[],"version":"current","sidebarPosition":2,"frontMatter":{"sidebar_position":2},"sidebar":"tutorialSidebar","previous":{"title":"Introduction","permalink":"/promptext/"},"next":{"title":"Configuration","permalink":"/promptext/configuration"}}');var l=s(4848),t=s(8453);const r={sidebar_position:2},o="Getting Started",a={},c=[{value:"Installation",id:"installation",level:2},{value:"Prerequisites",id:"prerequisites",level:3},{value:"All Platforms",id:"all-platforms",level:4},{value:"Platform-Specific",id:"platform-specific",level:4},{value:"Installation Methods",id:"installation-methods",level:3},{value:"Usage",id:"usage",level:2},{value:"Basic Command Structure",id:"basic-command-structure",level:3},{value:"Available Flags",id:"available-flags",level:3},{value:"Examples",id:"examples",level:3}];function d(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",h4:"h4",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,t.R)(),...e.components};return(0,l.jsxs)(l.Fragment,{children:[(0,l.jsx)(n.header,{children:(0,l.jsx)(n.h1,{id:"getting-started",children:"Getting Started"})}),"\n",(0,l.jsx)(n.h2,{id:"installation",children:"Installation"}),"\n",(0,l.jsx)(n.h3,{id:"prerequisites",children:"Prerequisites"}),"\n",(0,l.jsx)(n.h4,{id:"all-platforms",children:"All Platforms"}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"Git (for version control features)"}),"\n"]}),"\n",(0,l.jsx)(n.h4,{id:"platform-specific",children:"Platform-Specific"}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.strong,{children:"Linux/macOS"}),": No additional requirements"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.strong,{children:"Windows"}),":","\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"PowerShell 5.1 or higher"}),"\n"]}),"\n"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.strong,{children:"Go Installation"}),": Go 1.22 or higher (if installing via ",(0,l.jsx)(n.code,{children:"go install"}),")"]}),"\n"]}),"\n",(0,l.jsx)(n.h3,{id:"installation-methods",children:"Installation Methods"}),"\n",(0,l.jsxs)(n.ol,{children:["\n",(0,l.jsx)(n.li,{children:"Quick Install:"}),"\n"]}),"\n",(0,l.jsxs)(n.p,{children:[(0,l.jsx)(n.strong,{children:"Linux/macOS"}),":"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:"# User installation (recommended)\ncurl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash --user\n\n# System-wide installation (requires sudo)\ncurl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | sudo bash\n\n# Additional options\ncurl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash --help\n"})}),"\n",(0,l.jsx)(n.p,{children:"The Linux/macOS installer provides:"}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"User-level or system-wide installation"}),"\n",(0,l.jsx)(n.li,{children:"Custom installation directory support"}),"\n",(0,l.jsx)(n.li,{children:"Automatic checksum verification"}),"\n",(0,l.jsx)(n.li,{children:"Shell-specific alias configuration"}),"\n",(0,l.jsx)(n.li,{children:"Clean uninstallation"}),"\n",(0,l.jsx)(n.li,{children:"HTTPS security options"}),"\n"]}),"\n",(0,l.jsxs)(n.p,{children:[(0,l.jsx)(n.strong,{children:"Windows (PowerShell)"}),":"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-powershell",children:"# Install\nirm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex\n\n# Uninstall\nirm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex -Uninstall\n"})}),"\n",(0,l.jsx)(n.p,{children:"The installers provide:"}),"\n",(0,l.jsx)(n.p,{children:"Linux/macOS:"}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"\u2728 User-level or system-wide installation"}),"\n",(0,l.jsx)(n.li,{children:"\ud83d\udd12 Automatic checksum verification"}),"\n",(0,l.jsx)(n.li,{children:"\ud83d\udcc1 Custom installation directory support"}),"\n",(0,l.jsx)(n.li,{children:"\u26a1 PATH environment configuration"}),"\n",(0,l.jsx)(n.li,{children:"\ud83d\udcab Shell alias configuration"}),"\n",(0,l.jsx)(n.li,{children:"\ud83d\udd04 Clean uninstallation"}),"\n"]}),"\n",(0,l.jsx)(n.p,{children:"Windows:"}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"\ud83c\udfe0 User-level installation in %LOCALAPPDATA%"}),"\n",(0,l.jsx)(n.li,{children:"\ud83d\udd12 Automatic checksum verification"}),"\n",(0,l.jsx)(n.li,{children:"\u26a1 User PATH configuration"}),"\n",(0,l.jsx)(n.li,{children:"\ud83d\udcab PowerShell alias creation (prx)"}),"\n",(0,l.jsx)(n.li,{children:"\ud83d\udd04 Clean uninstallation"}),"\n"]}),"\n",(0,l.jsxs)(n.ol,{start:"2",children:["\n",(0,l.jsx)(n.li,{children:"Using Go Install:"}),"\n"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:"go install github.com/1broseidon/promptext/cmd/promptext@latest\n"})}),"\n",(0,l.jsxs)(n.ol,{start:"3",children:["\n",(0,l.jsx)(n.li,{children:"Manual Installation:"}),"\n"]}),"\n",(0,l.jsxs)(n.p,{children:["Download the appropriate binary for your platform from the ",(0,l.jsx)(n.a,{href:"https://github.com/1broseidon/promptext/releases",children:"releases page"}),":"]}),"\n",(0,l.jsxs)(n.p,{children:[(0,l.jsx)(n.strong,{children:"Linux/macOS"}),":"]}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"Download the appropriate binary"}),"\n",(0,l.jsxs)(n.li,{children:["Make it executable: ",(0,l.jsx)(n.code,{children:"chmod +x promptext"})]}),"\n",(0,l.jsxs)(n.li,{children:["Move to PATH: ",(0,l.jsx)(n.code,{children:"sudo mv promptext /usr/local/bin/"})]}),"\n"]}),"\n",(0,l.jsxs)(n.p,{children:[(0,l.jsx)(n.strong,{children:"Windows"}),":"]}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"Download the Windows binary (ZIP file)"}),"\n",(0,l.jsxs)(n.li,{children:["Extract to a directory (e.g., ",(0,l.jsx)(n.code,{children:"%LOCALAPPDATA%\\promptext"}),")"]}),"\n",(0,l.jsxs)(n.li,{children:["Add the directory to your PATH:","\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsx)(n.li,{children:"System Settings > Advanced > Environment Variables"}),"\n",(0,l.jsx)(n.li,{children:"Edit the User PATH variable"}),"\n",(0,l.jsx)(n.li,{children:"Add the installation directory"}),"\n"]}),"\n"]}),"\n"]}),"\n",(0,l.jsx)(n.h2,{id:"usage",children:"Usage"}),"\n",(0,l.jsx)(n.h3,{id:"basic-command-structure",children:"Basic Command Structure"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:"promptext [flags]\n"})}),"\n",(0,l.jsx)(n.h3,{id:"available-flags",children:"Available Flags"}),"\n",(0,l.jsxs)(n.ul,{children:["\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-version, -v"}),": Show version information"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-directory, -d string"}),': Directory to process (default ".")']}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-extension, -e string"}),': File extensions to include (comma-separated, e.g., ".go,.js")']}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-exclude, -x string"}),": Patterns to exclude (comma-separated)"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-format, -f string"}),": Output format (markdown/xml)"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-output, -o string"}),": Output file path"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-info, -i"}),": Show only project summary with token counts"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-verbose, -V"}),": Show full file contents"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-debug, -D"}),": Enable debug logging"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-gitignore, -g"}),": Use .gitignore patterns (default true)"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-use-default-rules, -u"}),": Use default filtering rules (default true)"]}),"\n",(0,l.jsxs)(n.li,{children:[(0,l.jsx)(n.code,{children:"-help, -h"}),": Show help message"]}),"\n"]}),"\n",(0,l.jsx)(n.h3,{id:"examples",children:"Examples"}),"\n",(0,l.jsxs)(n.ol,{children:["\n",(0,l.jsx)(n.li,{children:"Process specific file types:"}),"\n"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:"promptext -extension .go,.js\n"})}),"\n",(0,l.jsxs)(n.ol,{start:"2",children:["\n",(0,l.jsx)(n.li,{children:"Export as XML with debug info:"}),"\n"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:"promptext -format xml -output project.xml -debug\n"})}),"\n",(0,l.jsxs)(n.ol,{start:"3",children:["\n",(0,l.jsx)(n.li,{children:"Show project overview with token counts:"}),"\n"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:"promptext -info\n"})}),"\n",(0,l.jsxs)(n.ol,{start:"4",children:["\n",(0,l.jsx)(n.li,{children:"Process with exclusions:"}),"\n"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:'promptext -exclude "test/,vendor/" -V\n'})}),"\n",(0,l.jsxs)(n.ol,{start:"5",children:["\n",(0,l.jsx)(n.li,{children:"Check version:"}),"\n"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:"promptext -v  # Show version information\npromptext --version  # Same as above\n\n# Example output:\n# promptext version v0.2.4 (2024-12-19)\n"})}),"\n",(0,l.jsxs)(n.ol,{start:"5",children:["\n",(0,l.jsx)(n.li,{children:"Process all files including dependencies:"}),"\n"]}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-bash",children:'promptext -u=false -exclude "test/" # Disable default rules but keep test/ excluded\n'})})]})}function h(e={}){const{wrapper:n}={...(0,t.R)(),...e.components};return n?(0,l.jsx)(n,{...e,children:(0,l.jsx)(d,{...e})}):d(e)}},8453:(e,n,s)=>{s.d(n,{R:()=>r,x:()=>o});var i=s(6540);const l={},t=i.createContext(l);function r(e){const n=i.useContext(t);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(l):e.components||l:r(e.components),i.createElement(t.Provider,{value:n},e.children)}}}]);