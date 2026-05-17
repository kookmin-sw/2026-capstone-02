%% betterposterv2.cls
%% Copyright 2026 Max Stabel
%
% This work may be distributed and/or modified under the
% conditions of the LaTeX Project Public License, either version 1.3
% of this license or (at your option) any later version.
% The latest version of this license is in
%   https://www.latex-project.org/lppl.txt
% and version 1.3c or later is part of all distributions of LaTeX
% version 2008 or later.
%
% This work has the LPPL maintenance status `maintained'.
% 
% The Current Maintainer of this work is Max Stabel.
%
% This work consists of the files betterposterv2.cls and poster.tex
% Date: 2026-01-21

\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{betterposterv2}[2026/01/21 Better Poster v2 Class]

% --- 1. Base Class & Engine ---
% Define class options
\RequirePackage{kvoptions}
\RequirePackage{etoolbox}  % Programming facilities (needed for \ifdefempty)
\RequirePackage{iftex}
\SetupKeyvalOptions{
    family=bp,
    prefix=bp@
}
% Color scheme option (default, green, blue, red, purple)
\DeclareStringOption[default]{colorscheme}
% Debug option (shows layout frames and enables draft mode)
\DeclareBoolOption{debug}
% Font size option (default 42pt for A0 poster)
\DeclareStringOption[42pt]{fontsize}
% Paper size option (default a0paper)
% \DeclareStringOption[a0paper]{papersize}
% Orientation option (portrait or landscape)
\DeclareStringOption[portrait]{orientation}
% Pass unknown options to base class
\DeclareDefaultOption{%
    \ifx\CurrentOptionValue\relax
        \PackageWarningNoLine{\@currname}{%
            Unknown option '\CurrentOption'\MessageBreak
            is passed to base class 'scrartcl'%
        }%
        % Pass through any option to base class
        \expandafter\PassOptionsToClass{\expandafter\CurrentOption}{scrartcl}%
    \fi
}
\ProcessKeyvalOptions*

% Pass papersize to base class
% \PassOptionsToClass{\bp@papersize}{scrartcl}
% Pass orientation to base class if specified
\PassOptionsToClass{\bp@orientation}{scrartcl}
% We use scrartcl (KOMA-Script article) as the base.
\LoadClass[fontsize=\bp@fontsize, parskip=half]{scrartcl}

% Check for LuaLaTeX/XeLaTeX for modern font support
\ifboolexpr{not bool{luatex} and not bool{xetex}}{
    \ClassWarningNoLine{betterposterv2}{%
        You are not using LuaLaTeX or XeLaTeX.\MessageBreak
        Modern font features will be limited.\MessageBreak
        Recommended: Compile with LuaLaTeX for best results%
    }%
}{}

% --- Debug Mode ---
\ifbp@debug
    \RequirePackage{showframe}  % Shows boxes around text area, margin, header and footer
    \KOMAoptions{draft=true}    % Print markers for overfull boxes/lines
    \PackageInfo{betterposterv2}{Debug mode enabled: showframe loaded, draft mode active}
\fi

% --- 2. Essential Layout Packages ---
\RequirePackage[export]{adjustbox}% Extends graphicx with alignment options
\RequirePackage{geometry}         % For margins
\RequirePackage{xcolor}           % For colors
\RequirePackage{graphicx}         % For images
\RequirePackage{scrlayer-scrpage} % KOMA-Script headers/footers
\RequirePackage{caption}          % For \captionof command
\RequirePackage{siunitx}          % For SI units
\RequirePackage{import}		      % Relative path imports
\RequirePackage{xparse}           % For modern command definitions
\RequirePackage{xstring}          % For \IfEqCase command
\RequirePackage{xifthen}          % For conditional commands
\RequirePackage{booktabs}         % For better tables
\RequirePackage{array}            % For better arrays/matrices
\RequirePackage{multirow}         % For multirow table cells
\RequirePackage{multicol}          % For multi-column text
\RequirePackage{longtable}        % For multi-page tables
\RequirePackage{csquotes}         % For proper quotation marks
\RequirePackage[tracking=true]{microtype}% Improves text spacing/appearance
\microtypesetup{protrusion=true,expansion=true}% Character protrusion and font expansion
\captionsetup[figure]{hypcap=false, labelformat=empty}% Disable hyperref link to figure caption and remove label
\graphicspath{{./graphics}{./plots}{./images}{./img}{./schematics}{./pictures}{./imgs}{./figures}}

% --- 3. Typography & Fonts ---
% Font Loading: Fira family (Sans, Math, Code) - Modern, Open Source
\ifboolexpr{bool{luatex} or bool{xetex}}{
    \RequirePackage{fontspec}
    \RequirePackage{unicode-math}% Load after fontspec, before font declarations

    % Try loading Fira Sans
    \IfFontExistsTF{Fira Sans}{
        \setsansfont{Fira Sans}[
            UprightFont=*-Regular,
            BoldFont=*-Bold,
            ItalicFont=*-Italic,
            BoldItalicFont=*-BoldItalic,
            Scale=1.0
        ]
        % Use Fira Sans as the main font for the poster
        \setmainfont{Fira Sans}[
            UprightFont=*-Regular,
            BoldFont=*-Bold,
            ItalicFont=*-Italic,
            BoldItalicFont=*-BoldItalic,
            Scale=1.0
        ]
        % Define additional font series for more weight options
        \newfontfamily\firasansmedium{Fira Sans}[
            UprightFont=*-Medium,
            ItalicFont=*-MediumItalic,
            Scale=1.0
        ]
        \newfontfamily\firasansextrabold{Fira Sans}[
            UprightFont=*-ExtraBold,
            ItalicFont=*-ExtraBoldItalic,
            Scale=1.0
        ]
    }{
        % Fallback if Fira Sans is not installed
        \setsansfont{Arial}
        \setmainfont{Arial}
    }

    % Monospace: Try Fira Code, then Fira Mono, then fallback
    \IfFontExistsTF{Fira Code}{
        \setmonofont{Fira Code}[
            Scale=0.9,
            Contextuals=Alternate
        ]
    }{
        \IfFontExistsTF{Fira Mono}{
            \setmonofont{Fira Mono}[Scale=0.9]
        }{
            \setmonofont{DejaVu Sans Mono}[Scale=0.9]
        }
    }

    % Math font: Fira Math pairs perfectly with Fira Sans
    \IfFontExistsTF{Fira Math}{
        \setmathfont{Fira Math}[
            Scale=1.0,
            math-style=ISO,
            bold-style=ISO
        ]
    }{
        % Fallback to a sans-serif math font if Fira Math not available
        \setmathfont{Latin Modern Math}
    }

    % Dedicated font family for siunitx numbers to ensure consistent appearance
    \newfontfamily{\unitnumberfont}{Fira Sans}[
        UprightFont=*-Regular,
        Numbers=Lining,
        Scale=1.0
    ]
}{
    % PDFLaTeX fallback
    \RequirePackage[T1]{fontenc}
    \RequirePackage{helvet}
    \renewcommand{\familydefault}{\sfdefault}
}

% Configure siunitx for poster presentation
\sisetup{
    detect-all,                    % Detect font weight/shape
    mode=text,                     % Use text fonts not math
    text-family-to-math=true,      % Ensure consistency
    text-series-to-math=true,
    per-mode=symbol,               % Use / for division
    group-separator={,},           % Thousands separator
    group-minimum-digits=4,        % Group numbers ≥ 10000
}

% --- 4. Colors ---
% Define base color palette
\definecolor{bpWhite}{HTML}{FFFFFF}      % White

% Grey scale (ETH grey)
\definecolor{bpLightGrey}{HTML}{F1F1F1}   % ETH grey 10%
\definecolor{bpMediumGrey}{HTML}{E2E2E2}  % ETH grey 20%
\definecolor{bpDark}{HTML}{575757}        % ETH grey 120%

% Green (ETH green)
\definecolor{bpLightGreen}{HTML}{EEF1E7}   % ETH green 10%
\definecolor{bpMediumGreen}{HTML}{E0E3D0}  % ETH green 20%
\definecolor{bpDarkGreen}{HTML}{365213}    % ETH green 120%

% Blue (ETH blue)
\definecolor{bpLightBlue}{HTML}{E9EFF7}    % ETH blue 10%
\definecolor{bpMediumBlue}{HTML}{D3DEEF}   % ETH blue 20%
\definecolor{bpDarkBlue}{HTML}{08407E}     % ETH blue 120%

% Red (ETH red)
\definecolor{bpLightRed}{HTML}{F8EBEA}     % ETH red 10%
\definecolor{bpMediumRed}{HTML}{F1D7D5}    % ETH red 20%
\definecolor{bpDarkRed}{HTML}{96272D}      % ETH red 120%

% Purple (ETH purple)
\definecolor{bpLightPurple}{HTML}{F8E8F3}  % ETH purple 10%
\definecolor{bpMediumPurple}{HTML}{EFD0E3} % ETH purple 20%
\definecolor{bpDarkPurple}{HTML}{8C0A59}   % ETH purple 120%

% Functional color assignments (these are what get used in the template)
\colorlet{bpPageBackground}{bpLightGrey}
\colorlet{bpHeaderBackground}{bpDark}
\colorlet{bpFooterBackground}{bpDark}
\colorlet{bpLightBoxBackground}{bpWhite}
\colorlet{bpDarkBoxBackground}{bpMediumGrey}
\colorlet{bpHeaderText}{white}
\colorlet{bpFooterText}{white}

% Apply page background color
\pagecolor{bpPageBackground}

% Color scheme presets
\NewDocumentCommand{\SetColorScheme}{m}{%
    \IfEqCase{#1}{%
        {default}{%
                \colorlet{bpPageBackground}{bpLightGrey}%
                \colorlet{bpHeaderBackground}{bpDark}%
                \colorlet{bpFooterBackground}{bpDark}%
                \colorlet{bpLightBoxBackground}{bpWhite}%
                \colorlet{bpDarkBoxBackground}{bpMediumGrey}%
                \colorlet{bpHeaderText}{white}%
                \colorlet{bpFooterText}{white}%
            }%
            {green}{%
                \colorlet{bpPageBackground}{bpLightGreen}%
                \colorlet{bpHeaderBackground}{bpDarkGreen}%
                \colorlet{bpFooterBackground}{bpDarkGreen}%
                \colorlet{bpLightBoxBackground}{bpWhite}%
                \colorlet{bpDarkBoxBackground}{bpMediumGreen}%
                \colorlet{bpHeaderText}{white}%
                \colorlet{bpFooterText}{white}%
            }%
            {blue}{%
                \colorlet{bpPageBackground}{bpLightBlue}%
                \colorlet{bpHeaderBackground}{bpDarkBlue}%
                \colorlet{bpFooterBackground}{bpDarkBlue}%
                \colorlet{bpLightBoxBackground}{bpWhite}%
                \colorlet{bpDarkBoxBackground}{bpMediumBlue}%
                \colorlet{bpHeaderText}{white}%
                \colorlet{bpFooterText}{white}%
            }%
            {red}{%
                \colorlet{bpPageBackground}{bpLightRed}%
                \colorlet{bpHeaderBackground}{bpDarkRed}%
                \colorlet{bpFooterBackground}{bpDarkRed}%
                \colorlet{bpLightBoxBackground}{bpWhite}%
                \colorlet{bpDarkBoxBackground}{bpMediumRed}%
                \colorlet{bpHeaderText}{white}%
                \colorlet{bpFooterText}{white}%
            }%
            {purple}{%
                \colorlet{bpPageBackground}{bpLightPurple}%
                \colorlet{bpHeaderBackground}{bpDarkPurple}%
                \colorlet{bpFooterBackground}{bpDarkPurple}%
                \colorlet{bpLightBoxBackground}{bpWhite}%
                \colorlet{bpDarkBoxBackground}{bpMediumPurple}%
                \colorlet{bpHeaderText}{white}%
                \colorlet{bpFooterText}{white}%
            }%
    }[\PackageWarning{betterposterv2}{Unknown color scheme '#1'. Available: default, green, blue, red, purple}]%
    \pagecolor{bpPageBackground}%
}

% Individual color override commands
\NewDocumentCommand{\SetPageColor}{m}{\colorlet{bpPageBackground}{#1}\pagecolor{bpPageBackground}}
\NewDocumentCommand{\SetHeaderColor}{m}{\colorlet{bpHeaderBackground}{#1}}
\NewDocumentCommand{\SetFooterColor}{m}{\colorlet{bpFooterBackground}{#1}}
\NewDocumentCommand{\SetLightBoxColor}{m}{\colorlet{bpLightBoxBackground}{#1}}
\NewDocumentCommand{\SetDarkBoxColor}{m}{\colorlet{bpDarkBoxBackground}{#1}}
\NewDocumentCommand{\SetHeaderTextColor}{m}{\colorlet{bpHeaderText}{#1}}
\NewDocumentCommand{\SetFooterTextColor}{m}{\colorlet{bpFooterText}{#1}}

% Apply color scheme from class option
\AtEndPreamble{\SetColorScheme{\bp@colorscheme}}

% --- 5. Layout Lengths ---
\newlength{\bpHeaderHeight}
\newlength{\bpFooterHeight}
\newlength{\bpSideMargin}

\setlength{\bpHeaderHeight}{20cm}
\setlength{\bpFooterHeight}{10cm}
\setlength{\bpSideMargin}{3cm}

\NewDocumentCommand{\ApplyPosterGeometry}{}{%
    \geometry{
        \bp@orientation,
        top=\dimexpr\bpHeaderHeight+\baselineskip\relax,
        bottom=\dimexpr\bpFooterHeight+\baselineskip\relax,
        left=\bpSideMargin,
        right=\bpSideMargin,
        headheight=\bpHeaderHeight,
        headsep=\baselineskip,
        footskip=\baselineskip,
        marginparwidth=0pt
    }%
}
\ApplyPosterGeometry

\NewDocumentCommand{\SetHeaderHeight}{m}{\setlength{\bpHeaderHeight}{#1}\ApplyPosterGeometry}
\NewDocumentCommand{\SetFooterHeight}{m}{\setlength{\bpFooterHeight}{#1}\ApplyPosterGeometry}
\NewDocumentCommand{\SetSideMargin}{m}{\setlength{\bpSideMargin}{#1}\ApplyPosterGeometry}

% --- Font Size Variables ---
\newcommand{\bpMainFindingSize}{\Huge}
\newcommand{\bpTitleSize}{\Large}
\newcommand{\bpBoxTitleSize}{\large}
\newcommand{\bpSectionSize}{\Large}

\NewDocumentCommand{\SetMainFindingSize}{m}{\renewcommand{\bpMainFindingSize}{#1}}
\NewDocumentCommand{\SetTitleSize}{m}{\renewcommand{\bpTitleSize}{#1}}
\NewDocumentCommand{\SetBoxTitleSize}{m}{\renewcommand{\bpBoxTitleSize}{#1}}
\NewDocumentCommand{\SetSectionSize}{m}{\renewcommand{\bpSectionSize}{#1}}

% --- 6. Content Boxes (tcolorbox) ---
\RequirePackage{tcolorbox}
\tcbuselibrary{skins, raster, breakable}

\tcbset{
    bpbase/.style={
            enhanced,
            boxrule=0pt,
            left=1cm, right=1cm, top=1cm, bottom=1cm,
            toptitle=1cm, bottomtitle=0cm,
            coltitle=black,
            % arc=4mm,
            fonttitle=\bfseries\bpBoxTitleSize,
            fontupper=\normalsize,
            left skip=0pt,
            right skip=0pt,
            before skip=0.5\baselineskip,
            after skip=0.5\baselineskip
        },
    bpcolor/.is choice,
    bpcolor/grey/.style={
            colback=bpDarkBoxBackground,
            colframe=bpDarkBoxBackground,
        },
    bpcolor/white/.style={
            colback=bpLightBoxBackground,
            colframe=bpLightBoxBackground,
        },
    bpcolor/.default=white,
    bptitlebg/.is choice,
    bptitlebg/normal/.style={},
    bptitlebg/transparent/.style={
            colbacktitle=bpPageBackground,
            bottomtitle=0.5cm,
        },
    bptitlebg/.default=normal,
    bpstyle/.is choice,
    bpstyle/normal/.style={},
    bpstyle/wide/.style={
            width=\linewidth,
        },
    bpstyle/.default=normal,
    bpbox/.style={
            bpbase,
            bpcolor=white,
            bpstyle=normal,
        }
}

\newtcolorbox{bpbox}[1][]{
    bpbox,
    #1
}

% --- 7. Grid Layout ---
\NewDocumentEnvironment{bpcolumns}{O{2}}{%
    % \par\noindent%
    \begin{tcbraster}[%
            raster columns=#1,%
            raster equal height=rows,%
            raster column skip=0.5\baselineskip,%
            raster row skip=0pt,%
            raster left skip=0pt,%
            raster right skip=0pt,%
            raster width=\linewidth,%
            raster before skip=0.5\baselineskip,%
            raster after skip=0.5\baselineskip,%
            raster valign=top%
        ]%
        }{%
    \end{tcbraster}
}

% --- 8. Section Header ---
\NewDocumentCommand{\bpsection}{m}{%
    \par\noindent\hspace{1cm}{\bpSectionSize\bfseries\color{black} #1}\par%
}

% --- 9. Header & Footer Content ---
\NewDocumentCommand{\mainfinding}{m}{\gdef\bpMainFinding{#1}}
\NewDocumentCommand{\bpMainFinding}{}{Insert Main Finding Here}

% Footer element commands
\NewDocumentCommand{\qrcode}{m}{\gdef\bpQRCode{#1}}
\newcommand{\bpQRCode}{}

\NewDocumentCommand{\qrcodelabel}{m}{\gdef\bpQRCodeLabel{#1}}
\newcommand{\bpQRCodeLabel}{}

\NewDocumentCommand{\authors}{m}{\gdef\bpAuthors{#1}}
\newcommand{\bpAuthors}{}

\NewDocumentCommand{\affiliations}{+m}{\gdef\bpAffiliations{#1}}
\newcommand{\bpAffiliations}{}

\NewDocumentCommand{\portrait}{m}{\gdef\bpPortrait{#1}}
\newcommand{\bpPortrait}{}

\NewDocumentCommand{\publicationyear}{m}{\gdef\bpYear{#1}}
\newcommand{\bpYear}{}

\NewDocumentCommand{\logo}{m}{\gdef\bpLogo{#1}}
\newcommand{\bpLogo}{}

\NewDocumentCommand{\footer}{m}{\gdef\bpFooterContent{#1}}
\newcommand{\bpFooterLayout}{%
    \parbox[c][\bpFooterHeight][c]{\dimexpr\textwidth+\bpSideMargin\relax}{%  % Go up to the edge on the right
        \color{bpFooterText}%
        \centering
        % Left-aligned content: QR code + Authors/Affiliations (natural widths)
        \ifdefempty{\bpQRCode}{}{%
            \begin{tabular}[c]{@{}c@{}}
                \includegraphics[height=0.8\bpFooterHeight, keepaspectratio]{\bpQRCode}%
                \ifdefempty{\bpQRCodeLabel}{}{ \\[-0.7em]{\tiny\bpQRCodeLabel}}%
            \end{tabular}%
            \hspace{1cm}%
        }%
        %
        % Authors & Affiliations - vertically centered, no fixed width
        {\large\begin{tabular}[c]{@{}l@{}}
                \ifdefempty{\bpAuthors}{}{\bpAuthors}%
                \ifdefempty{\bpAffiliations}{}{\ifdefempty{\bpAuthors}{}{\\}{\normalsize\bpAffiliations}}%
                \ifdefempty{\bpYear}{}{\ifdefempty{\bpAuthors}{\ifdefempty{\bpAffiliations}{}{\\}}{\\}{\normalsize\bpYear}}%
            \end{tabular}}%
        %
        \hfill%
        %
        % Right-aligned: Logo & Portrait (natural widths)
        \ifdefempty{\bpLogo}{}{%
            \includegraphics[height=\bpFooterHeight, keepaspectratio, valign=c]{\bpLogo}%
        }%
        \ifdefempty{\bpPortrait}{}{%
            \includegraphics[height=\bpFooterHeight, keepaspectratio, valign=c]{\bpPortrait}%
        }%
    }%
}

% --- 10. Hyperlinks & PDF Metadata ---
\RequirePackage[pdfusetitle]{hyperref}
\hypersetup{%
    hidelinks,
    plainpages=false,
    pdfkeywords={poster, scientific, research},
    pdfcreator={LuaLaTeX with betterposterv2 class},
    pdfdisplaydoctitle=true,
    pdfpagelayout=SinglePage,
    bookmarksnumbered=false,
}%

% Set PDF author and subject at begin document when macros are defined
\AtBeginDocument{%
    \hypersetup{%
        pdfauthor={\bpAuthors},
        pdfsubject={\bpMainFinding}
    }%
}

% --- 11. Header & Footer Layers ---
\clearpairofpagestyles  % clear existing header/footer, e.g. page numbering

% --- Backgrounds Layers ---
\DeclareNewLayer[
    background,
    page,
    height=\bpHeaderHeight,
    contents={\color{bpHeaderBackground}\rule{\paperwidth}{\bpHeaderHeight}}
]{bp.header.bg}

\DeclareNewLayer[
    background,
    page,
    voffset=\dimexpr\paperheight-\bpFooterHeight\relax,
    height=\bpFooterHeight,
    contents={\color{bpFooterBackground}\rule{\paperwidth}{\bpFooterHeight}}
]{bp.footer.bg}

% --- Content Layers ---
\DeclareNewLayer[
    foreground,
    head,
    height=\bpHeaderHeight,
    contents={%
            \parbox[c][\bpHeaderHeight][c]{\textwidth}{%
                \color{bpHeaderText}%
                \bpMainFindingSize%
                \bpMainFinding%
            }%
        }
]{bp.header.content}

\DeclareNewLayer[
    foreground,
    page,
    hoffset=\bpSideMargin,
    voffset=\dimexpr\paperheight-\bpFooterHeight\relax,
    width=\dimexpr\paperwidth-\bpSideMargin\relax,   % go up to the page edge on the right
    height=\bpFooterHeight,
    contents={\bpFooterLayout}
]{bp.footer.content}

\AddLayersToPageStyle{scrheadings}{
    bp.header.bg,
    bp.header.content,
    bp.footer.bg,
    bp.footer.content
}
\pagestyle{scrheadings}

% --- 12. Title Block ---
\RenewDocumentCommand{\maketitle}{}{%
    \thispagestyle{scrheadings}%
    \bpTitleSize%
    \@title
}

\endinput