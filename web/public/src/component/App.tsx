import React from "react";
import { GetList } from "./repository/Novel";

interface AppProps {
  id: number;
  ids: number[];
}
const App = (props: AppProps) => {
  const a = GetList();
  return (
    <div className={"reactApp"}>
      HELLO React4
      {props.id}
      <br />
    </div>
  );
};
export default App;
