
import "./Header.css";
import "./Main.css";
import "./Footer.css";
import "./App.css";

import { useState, useEffect, useRef } from "react";

import { EditorView, basicSetup } from "codemirror";
import { EditorState, StateField, StateEffect } from "@codemirror/state";
import { Decoration, type DecorationSet } from "@codemirror/view";
import { go } from "@codemirror/lang-go";
import { dracula } from 'thememirror';

import mermaid from "mermaid";
import panzoom from "panzoom";

class Node_t {
    id: number;
    code: string;
    type: string;
    line: number;
    state: string;

    constructor(id: number, code: string, type: string, line: number, state: string = "") {
        this.id = id;
        this.code = code;
        this.type = type;
        this.line = line;
        this.state = state;
    }
};

class Edge_t {
    id: number;
    from: number;
    dest: number;
    cond: string;

    constructor(id: number, from: number, dest: number, cond: string) {
        this.id = id;
        this.from = from;
        this.dest = dest;
        this.cond = cond;
    }
};

class Func_t {
    name: string;
    nodes: Array<Node_t>;
    edges: Array<Edge_t>;

    constructor(name: string, nodes: Array<Node_t> = [], edges: Array<Edge_t> = []) {
        this.name = name;
        this.nodes = nodes;
        this.edges = edges;
    }
};

class Mermaid_t {
    functions: Array<Func_t>;

    constructor(functions: Array<Func_t>) {
        this.functions = functions;
    }
};

class Debug_t {
    type: string;
    funcName: string;
    nodeID: number;
    nodeState: string;
    message: string;
    lineNum: number;

    constructor(type: string, funcName: string, nodeID: number, nodeState: string, message: string, lineNum: number) {
        this.type = type;
        this.funcName = funcName;
        this.nodeID = nodeID;
        this.nodeState = nodeState;
        this.message = message;
        this.lineNum = lineNum;
    }
};

const highlightLineEffect = StateEffect.define<number>();
const highlightLineField = StateField.define<DecorationSet>({
    create() {
        return Decoration.none;
    },
    update(deco, tr) {
        deco = deco.map(tr.changes);

        for (let e of tr.effects) {
            if (e.is(highlightLineEffect)) {
                const line = e.value;

                const decoration = Decoration.line({
                    attributes: { class: "cm-highlight-line" }
                });

                const linePos = tr.state.doc.line(line);

                deco = Decoration.set([
                    decoration.range(linePos.from)
                ]);
            }
        }

        return deco;
    },
    provide: f => EditorView.decorations.from(f)
});

function App() {
    const codeEditorRef = useRef<HTMLDivElement>(null);
    const codeViewRef = useRef<EditorView>(null);

    const tabNames = useRef<string[]>([]);
    const [activeTab, setActiveTab] = useState<number>(-1);

    const mermaidRef = useRef<Mermaid_t>(null);
    const mermaidSrcRef = useRef<HTMLDivElement>(null);
    const [mermaidSrcs, setMermaidSrcs] = useState<string[]>([]);

    const [updateNodeSteps, setUpdateNodeSteps] = useState<Debug_t[]>([]);
    const [mermaidSrcsWithState, setMermaidSrcsWithState] = useState<string[]>([]);
    const [currentStepIndex, setCurrentStepIndex] = useState<number>(0);

    const fileRef = useRef<File>(null);
    const [_fileContent, setFileContent] = useState<string>("");
    const [_fileName, setFileName] = useState<string>("");

    const [leftWidth, setLeftWidth] = useState(33.3);
    const [middleWidth, setMiddleWidth] = useState(33.3);
    const [rightWidth, setRightWidth] = useState(33.3);

    const [isDragging, setIsDragging] = useState<null | "left" | "right">(null);

    const panzoomInstanceRef = useRef<ReturnType<typeof panzoom>>(null);

    // Create code editor view
    useEffect(() => {
        if (!codeEditorRef.current) return;

        const view = new EditorView({
            doc: "",
            parent: codeEditorRef.current,
            extensions: [
                basicSetup,
                highlightLineField,
                dracula,
                EditorState.readOnly.of(true),
                EditorView.editable.of(false),
                EditorView.contentAttributes.of({ tabindex: "0" }),
                go(),
            ]
        });

        codeViewRef.current = view;

        return () => {
            view.destroy();
        };
    }, []);

    // Initialize mermaid functionality
    useEffect(() => {
        mermaid.initialize({
            startOnLoad: false,
            securityLevel: 'strict',
            markdownAutoWrap: false,
            theme: "base",
            themeVariables: {
                primaryColor: '#44475a',
                primaryTextColor: '#f8f8f2',
                primaryBorderColor: '#6272a4',
                lineColor: '#6272a4',
                secondaryColor: '#282a36',
                tertiaryColor: '#44475a',
                background: '#282a36',
                mainBkg: '#282a36',
                secondBkg: '#44475a',
                tertiaryBkg: '#6272a4',
                edgeLabelBackground: '#282a36'
            }
        });
    }, []);

    // Call function to render flowchart by condition
    useEffect(() => {
        if (mermaidSrcs.length === 0) {
            if (mermaidSrcRef.current) {
                mermaidSrcRef.current.innerHTML = "";
            }
            return;
        }

        // If debug step mode is active, do not render normal charts
        if (updateNodeSteps.length > 0) {
            return;
        }

        if (activeTab >= 0 && activeTab < mermaidSrcs.length) {
            renderMermaid(activeTab, false);
        }
    }, [activeTab, mermaidSrcs, updateNodeSteps]);

    // Call function to render flowchart by another condition
    useEffect(() => {
        if (
            updateNodeSteps.length > 0 &&
            currentStepIndex >= 0 &&
            currentStepIndex < mermaidSrcsWithState.length
        ) {
            renderMermaid(currentStepIndex, true);
        }
    }, [currentStepIndex, mermaidSrcsWithState]);

    // Control slider with arrow keys
    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            if (updateNodeSteps.length === 0)
                return;

            if (event.key !== "ArrowLeft" && event.key !== "ArrowRight") {
                return;
            }

            event.preventDefault();

            let nextIndex = currentStepIndex;

            if (event.key === "ArrowRight") {
                nextIndex = Math.min(currentStepIndex + 1, updateNodeSteps.length - 1);
            }
            else if (event.key === "ArrowLeft") {
                nextIndex = Math.max(currentStepIndex - 1, 0);
            }

            if (nextIndex !== currentStepIndex) {
                handleSliderChange(nextIndex);
            }
        };

        window.addEventListener("keydown", handleKeyDown);

        return () => {
            window.removeEventListener("keydown", handleKeyDown);
        };
    }, [currentStepIndex, updateNodeSteps]);

    // Highlight line of the code background using timeline control
    useEffect(() => {
        if (!codeViewRef.current) return;

        const step = updateNodeSteps[currentStepIndex];
        if (!step) return;

        codeViewRef.current.dispatch({
            effects: highlightLineEffect.of(step.lineNum)
        });

    }, [currentStepIndex, updateNodeSteps]);

    // Control size of container boxes using mouse
    useEffect(() => {
        const handleMouseMove = (e: MouseEvent) => {
            if (!isDragging)
                return;

            const totalWidth = window.innerWidth;
            const mousePercent = (e.clientX / totalWidth) * 100;

            // Left divider
            if (isDragging === "left") {
                const newLeftWidth = Math.min(Math.max(mousePercent, 15), 70);
                const remaining = 100 - newLeftWidth - rightWidth;

                if (remaining >= 15) {
                    setLeftWidth(newLeftWidth);
                    setMiddleWidth(remaining);
                }
            }

            // Right divider
            else if (isDragging === "right") {
                const newRightWidth = Math.min(Math.max(100 - mousePercent, 15), 70);
                const remaining = 100 - leftWidth - newRightWidth;

                if (remaining >= 15) {
                    setRightWidth(newRightWidth);
                    setMiddleWidth(remaining);
                }
            }
        };

        const handleMouseUp = () => {
            setIsDragging(null);
        };

        window.addEventListener("mousemove", handleMouseMove);
        window.addEventListener("mouseup", handleMouseUp);

        return () => {
            window.removeEventListener("mousemove", handleMouseMove);
            window.removeEventListener("mouseup", handleMouseUp);
        };
    }, [isDragging, leftWidth, rightWidth]);

    // Trigger hidden file input
    const handleFileClick = () => {
        const fileInput = document.getElementById("fileInput") as HTMLInputElement;
        fileInput?.click();
    };

    // Read file content and print
    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];

        if (!file)
            return;

        setFileName(file.name);

        fileRef.current = file

        if (!file.name.toLowerCase().endsWith(".go"))
            return;

        const reader = new FileReader();
        reader.onload = (e) => {
            const content = e.target?.result as string;

            codeViewRef.current?.dispatch({
                changes: {
                    from: 0,
                    to: codeViewRef.current.state.doc.length,
                    insert: content
                }
            });

            setFileContent(content);
        };
        reader.readAsText(file);
    };

    // Upload file to server
    const handleFileUpload = async () => {
        if (!fileRef.current) {
            alert("No file selected!");
            return;
        }

        const form = new FormData();
        form.append("file", fileRef.current);

        try {
            const res = await fetch("/upload", {
                method: "POST",
                body: form
            });

            if (!res.ok) {
                const err = await res.text();
                console.error(`Upload failed!\n${err}`);
                alert(`Upload failed!\n${err}`);
                return;
            }

        } catch (err) {
            console.error(`Upload error!\n${err}`);
            alert(`Upload error\n${err}`);
        }
    };

    // Run inspector and print mermaids
    const handlePrintMermaids = async () => {
        if (!fileRef.current) {
            alert("No file selected!");
            return;
        }

        try {
            const res = await fetch("/run", {
                method: "GET"
            });

            if (!res.ok) {
                await res.text();
                console.error(`Flowchart generation failed!\n`);
                alert(`Flowchart generation failed!\n`);
                return;
            }

            const data = await res.json();
            const output = data.output;

            let convMermaidSrc: string = "";
            let convMermaidSrcs: Array<string> = [];
            let convFuncs: Array<Func_t> = [];

            tabNames.current = [];

            // Convert JSON to mermaid objects
            for (const outFuncsName in output) {
                const outFuncs = output[outFuncsName];

                tabNames.current.push(outFuncsName);

                convMermaidSrc += `flowchart TB\n`;

                let convNodes: Array<Node_t> = [];
                let convEdges: Array<Edge_t> = [];

                if (outFuncs.Nodes) {
                    // Clone and sort outNodes
                    const outNodes = [...outFuncs.Nodes].sort((a, b) => {
                        const ai = Number(a.Id);
                        const bi = Number(b.Id);

                        if (Number.isNaN(ai) || Number.isNaN(bi))
                            return String(a.Id).localeCompare(String(b.Id));

                        return ai - bi;
                    });

                    // Convert outNodes to node object
                    for (let i = 0; i < outNodes.length; i++) {
                        const outNodeID = outNodes[i].Id;
                        const outNodeCode = outNodes[i].Code;
                        const outSafeCode = outNodeCode.replaceAll("`", "#96;").replaceAll("\"", "#34;");
                        const outNodeType = outNodes[i].Node_type;
                        const outLineNum = outNodes[i].Line_num;

                        convMermaidSrc += `    id${outNodeID}`

                        if (outNodeType === "basic") {
                            convMermaidSrc += `[\"\`${outSafeCode}\`\"]`;
                        }
                        else if (outNodeType === "cond") {
                            convMermaidSrc += `{\"\`${outSafeCode}\`\"}`;
                        }

                        convMermaidSrc += `\n`;

                        convNodes.push(new Node_t(outNodeID, outSafeCode, outNodeType, outLineNum));
                    }
                }

                if (outFuncs.Edges) {
                    // Clone and sort outEdges
                    const outEdges = [...outFuncs.Edges].sort((a, b) => {
                        const ai = Number(a.Id);
                        const bi = Number(b.Id);

                        if (Number.isNaN(ai) || Number.isNaN(bi))
                            return String(a.Id).localeCompare(String(b.Id));

                        return ai - bi;
                    });

                    // Convert outEdges to edge object
                    for (let i = 0; i < outEdges.length; i++) {
                        const outEdgeID = outEdges[i].Id;
                        const outEdgeFrom = outEdges[i].From_node_id;
                        const outEdgeDest = outEdges[i].To_node_id;
                        const outEdgeCond = outEdges[i].Label;

                        convMermaidSrc += `    id${outEdgeFrom} `;
                        convMermaidSrc += outEdgeCond !== "" ? `-- ${outEdgeCond} --> ` : `--> `;
                        convMermaidSrc += `id${outEdgeDest}\n`;

                        convEdges.push(new Edge_t(outEdgeID, outEdgeFrom, outEdgeDest, outEdgeCond));
                    }
                }

                convMermaidSrcs.push(convMermaidSrc);
                convMermaidSrc = "";

                convFuncs.push(new Func_t(outFuncsName, convNodes, convEdges));
            }

            // Set mermaid objects to globally
            mermaidRef.current = new Mermaid_t(convFuncs);

            // Set mermaid source code to globally
            setMermaidSrcs(convMermaidSrcs);

            // Set default tab
            setActiveTab(convMermaidSrcs.length ? 0 : -1);

        } catch (_err) {
            console.error(`Flowchart generation error!\n`);
            alert(`Flowchart generation error!\n`);
        }
    };

    // Format node state to look clean
    const formatNodeState = (state: string): string => {
        if (!state || state === "{}")
            return state;

        let trimmed = state.trim();

        // Remove only outermost braces
        if (trimmed.startsWith("{") && trimmed.endsWith("}")) {
            trimmed = trimmed.slice(1, -1);
        }

        let result = "";
        let depth = 0;

        for (let i = 0; i < trimmed.length; i++) {
            const ch = trimmed[i];

            // Split only on top-level commas
            if (ch === "," && depth === 0) {
                result += "\n";
                continue;
            }

            result += ch;

            // Update nesting AFTER processing current char
            if (ch === "{" || ch === "[" || ch === "(") {
                depth++;
            }
            else if (ch === "}" || ch === "]" || ch === ")") {
                depth--;
            }
        }

        return result
            .split("\n")
            .map(line => line.trim())
            .join("\n");
    };

    // Run inspector again and print debugs
    const handlePrintDebugs = async (): Promise<Debug_t[]> => {
        if (!fileRef.current) {
            alert("No file selected!");
            return [];
        }

        try {
            const res = await fetch("/debug", {
                method: "GET"
            });

            if (!res.ok) {
                await res.text();
                console.error(`Debug output generation failed!\n`);
                alert(`Debug output generation failed!\n`);
                return [];
            }

            const data = await res.json();
            const output = data.output;

            let outDebugs: Array<Debug_t> = [];

            // Convert JSON to array of debug objects
            if (output.Debugs) {
                for (let i = 0; i < output.Debugs.length; i++) {
                    const outType: string = output.Debugs[i].Type;
                    const outFuncName: string = output.Debugs[i].Function_name;
                    const outNodeID: number = output.Debugs[i].Node_id;
                    const outNodeState: string = formatNodeState(output.Debugs[i].Node_state);
                    const outMessage: string = output.Debugs[i].Msg;
                    const outLineNum: number = output.Debugs[i].Line_number;

                    outDebugs.push(new Debug_t(outType, outFuncName, outNodeID, outNodeState, outMessage, outLineNum));
                }

                // Get HTML div element
                const logBox = document.getElementById("logBox");

                if (logBox) {
                    // Clear previous content
                    logBox.innerHTML = "";

                    // Create HTML table
                    const table = document.createElement("table");
                    table.className = "debug-table";

                    // Create thead
                    const thead = table.createTHead();
                    const headerRow = thead.insertRow();
                    ["Type", "Function", "Line", "Message"].forEach((h) => {
                        const th = document.createElement("th");
                        th.textContent = h;
                        headerRow.appendChild(th);
                    });

                    // Create tbody
                    const tbody = table.createTBody();

                    // Append rows
                    for (let i = 0; i < outDebugs.length; i++) {
                        if (outDebugs[i].type !== "update_node") {
                            const row = tbody.insertRow();
                            const cType = row.insertCell();
                            const cFuncName = row.insertCell();
                            const cLineNum = row.insertCell();
                            const cMessage = row.insertCell();

                            const isWarning = outDebugs[i].type === "warning";
                            const isError = outDebugs[i].type === "error";

                            const iconText = isWarning ? "\uea6c" : isError ? "\uea87" : "\uea74";
                            const iconSpan = document.createElement("span");
                            iconSpan.textContent = iconText;
                            iconSpan.className = `icon-${outDebugs[i].type}`;

                            cType.appendChild(iconSpan);
                            cFuncName.textContent = outDebugs[i].funcName;
                            cLineNum.textContent = outDebugs[i].lineNum.toString();
                            cMessage.textContent = outDebugs[i].message;

                            row.className = `debug-${outDebugs[i].type}`;
                        }
                    }

                    // Append table to logBox
                    logBox.appendChild(table);
                }
            }

            return outDebugs;

        } catch (_err) {
            console.error(`Debug output generation error!\n`);
            alert(`Debug output generation error!\n`);
            return [];
        }
    };

    // Configure Debug Changes
    const configureDebugChanges = async (debugData: Debug_t[]) => {
        // Extract "update_node" debug objects
        const updateNodes = debugData.filter(d => d.type === "update_node");
        setUpdateNodeSteps(updateNodes);

        // Generate lots of "update_node" mermaidJS source codes
        if (updateNodes.length > 0) {
            if (!mermaidRef.current) {
                return;
            }

            const srcsWithState: string[] = [];

            for (let i = 0; i < updateNodes.length; i++) {
                const currentFunc = mermaidRef.current.functions.find(
                    f => f.name === updateNodes[i].funcName
                );

                if (!currentFunc)
                    continue;

                const targetNodeID = updateNodes[i].nodeID;
                const targetNodeState = updateNodes[i].nodeState;

                let convMermaidSrc = "flowchart TB\n";

                for (const outNode of currentFunc.nodes) {
                    const outNodeID = outNode.id;
                    const outNodeCode = outNode.code;
                    const outNodeType = outNode.type;

                    convMermaidSrc += `    id${outNodeID}`

                    if (outNodeID === targetNodeID) {
                        if (targetNodeState !== "{}") {
                            outNode.state = targetNodeState;
                        }
                    }

                    const formattedState = outNode.state
                        ? `<span style='color:#f1fa8c'>${outNode.state.replaceAll("\n", "<br/>")}</span><br/><br/>`
                        : "";

                    if (outNodeType === "basic") {
                        convMermaidSrc += `[\"\`${formattedState}${outNodeCode}\`\"]`;
                    }
                    else if (outNodeType === "cond") {
                        convMermaidSrc += `{\"\`${formattedState}${outNodeCode}\`\"}`;
                    }

                    convMermaidSrc += `\n`;
                }

                for (const outEdge of currentFunc.edges) {
                    const outEdgeFrom = outEdge.from;
                    const outEdgeDest = outEdge.dest;
                    const outEdgeCond = outEdge.cond;

                    convMermaidSrc += `    id${outEdgeFrom} `;
                    convMermaidSrc += outEdgeCond !== "" ? `-- ${outEdgeCond} --> ` : `--> `;
                    convMermaidSrc += `id${outEdgeDest}\n`;
                }

                srcsWithState.push(convMermaidSrc);
            }

            setMermaidSrcsWithState(srcsWithState);

            setCurrentStepIndex(0);

            const firstStep = updateNodes[0];

            if (firstStep) {
                const firstTabIndex = tabNames.current.findIndex((name) => name === firstStep.funcName);

                if (firstTabIndex !== -1) {
                    setActiveTab(firstTabIndex);
                }
            }
        }
    };

    const handleSliderChange = (newIndex: number) => {
        setCurrentStepIndex(newIndex);

        const step = updateNodeSteps[newIndex];

        if (step) {
            const tabIndex = tabNames.current.findIndex((name) => name === step.funcName);

            if (tabIndex !== -1) {
                setActiveTab(tabIndex);
            }
        }
    };

    // Render flowchart for the active tab
    const renderMermaid = async (index: number, isStepView: boolean = false) => {
        if (!mermaidSrcRef.current) return;

        if (isStepView && (index < 0 || index >= mermaidSrcsWithState.length)) return;
        if (!isStepView && (index < 0 || index >= mermaidSrcs.length)) return;

        const mermaidId = `mermaid-${isStepView ? "step" : "tab"}-${index}`;

        const srcToRender = isStepView ? mermaidSrcsWithState[index] : mermaidSrcs[index];

        try {
            const { svg } = await mermaid.render(mermaidId, srcToRender);

            mermaidSrcRef.current.innerHTML = `<div id="${mermaidId}" class="mermaid-wrapper">${svg}</div>`;

            if (panzoomInstanceRef.current) {
                panzoomInstanceRef.current.dispose();
            }

            const svgElement = mermaidSrcRef.current.querySelector("svg");

            if (svgElement) {
                panzoomInstanceRef.current = panzoom(svgElement, {
                    maxZoom: 5,
                    minZoom: 0.5,

                    smoothScroll: false,

                    zoomSpeed: 0.08,
                    zoomDoubleClickSpeed: 1,

                    beforeMouseDown: (_e: MouseEvent) => false,
                    beforeWheel: (_e: WheelEvent) => false,

                    filterKey: () => true
                });
            }

            requestAnimationFrame(() => {
                const nodeGroups = mermaidSrcRef.current?.querySelectorAll("g[class*='node']");

                if (!nodeGroups) return;

                nodeGroups.forEach((nodeGroup: Element) => {
                    const gElement = nodeGroup as SVGGElement;
                    const nodeId = gElement.id?.replace(/^id/, "");

                    if (nodeId && !isNaN(Number(nodeId))) {
                        gElement.setAttribute("data-id", nodeId);
                        gElement.style.cursor = "pointer";
                    }
                });
            });

        } catch (err) {
            console.error(`Mermaid render error\n${err}`);
        }
    };

    // Run inspector and print
    const handleRunInspection = async () => {
        if (!fileRef.current) {
            alert("No file selected!");
            return;
        }

        try {
            // Upload file to server
            await handleFileUpload();

            // Run inspector and print mermaids
            await handlePrintMermaids();

            // Run inspector and print debugs
            const debugData = await handlePrintDebugs();

            // Configure Debug Changes
            await configureDebugChanges(debugData);
            
        } catch (_err) {
            console.error(`Inspection error!\n`);
            alert(`Inspection error!\n`);
        }
    };

    return (
        <>
            <header>
                <div className="headBoxes" id="headTitleBox">
                    <div id="headTitle">TraceInspector</div>
                </div>
                <div className="headBoxes" id="headButtonBox">
                    <button onClick={handleFileClick} className="headButtons" id="openButton">Open</button>
                    <input id="fileInput" type="file" accept=".go" onChange={handleFileChange} style={{ display: "none" }} />
                    <button className="headButtons" id="runButton" onClick={handleRunInspection}>Run</button>
                </div>
                <br /><br /><br /><br /><hr /><br />
            </header>
            <main className="main-layout">
                <div ref={codeEditorRef} className="mainBoxes" id="codeBox" style={{ width: `${leftWidth}%` }} />
                <div className="resize-bar" onMouseDown={() => setIsDragging("left")} />
                <div className="mainBoxes" id="mermaidBox" style={{ width: `${middleWidth}%` }}>
                    {mermaidSrcs.length > 1 && (
                        <div className="mermaid-tabs">
                            {mermaidSrcs.map((_callback, index) => (
                                <button
                                    key={index}
                                    className={index === activeTab ? "tab active" : "tab"}
                                    onClick={() => {
                                        setActiveTab(index);

                                        const firstStepIndex = updateNodeSteps.findIndex(
                                            (step) => step.funcName === tabNames.current[index]
                                        );

                                        if (firstStepIndex !== -1) {
                                            setCurrentStepIndex(firstStepIndex);
                                        }
                                    }}
                                >
                                    {`${tabNames.current[index]}`}
                                </button>
                            ))}
                        </div>
                    )}
                    <div ref={mermaidSrcRef} className="mermaid-container" />
                </div>
                <div className="resize-bar" onMouseDown={() => setIsDragging("right")} />
                <div className="mainBoxes" id="logBox" style={{ width: `${rightWidth}%` }} />
            </main>
            <footer>
                {updateNodeSteps.length > 0 && (
                    <div className="step-control">
                        <div className="step-attr">
                            {currentStepIndex + 1} / {updateNodeSteps.length} | 
                            Functions: {updateNodeSteps[currentStepIndex]?.funcName} |
                            Line: {updateNodeSteps[currentStepIndex]?.lineNum} |
                            Node State: {(updateNodeSteps[currentStepIndex]?.nodeState !== "{}") ?
                                <span style={{ color: `#f1fa8c` }}>{updateNodeSteps[currentStepIndex]?.nodeState}</span> :
                                <span style={{ color: `#f1fa8c` }}>None</span>}
                        </div>
                        <br /><br /><br /><br />
                        <input
                            type="range"
                            min="0"
                            max={updateNodeSteps.length - 1}
                            value={currentStepIndex}
                            onChange={(e) => {
                                const newIndex = Number(e.target.value);
                                handleSliderChange(newIndex);
                            }}
                            className="step-slider"
                        />
                    </div>
                )}
            </footer>
        </>
    )
}

export default App;
