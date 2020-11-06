import React from "react";

interface AppProps {
  id: number;
  ids: number[];
}
const App = (props: AppProps) => {
  return (
    <div className={"reactApp"}>
      HELLO React4
      {props.id}
      <br />
      {props.ids}
    </div>
  );
};
export default App;
const grpcAddress = "localhost:18080";
