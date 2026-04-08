
import "./Header.css"

function Header() {
    return (
        <header>
            <div className="headBoxes" id="headTitleBox">
                <div id="headTitle">TraceInspector</div>
            </div>
            <div className="headBoxes" id="headButtonBox">
                <button className="headButtons" id="openButton">Open</button>
                <button className="headButtons" id="runButton">Run</button>
            </div>
            <br /><br /><br /><br /><br /><br /><hr /><br />
        </header>
    );
}

export default Header;
