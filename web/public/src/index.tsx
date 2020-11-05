import React from "react";
import ReactDOM from "react-dom";
import App from './component/App';


ReactDOM.render(<App id={1} ids={[1,2,3]}/>,
    document.getElementById('app'))