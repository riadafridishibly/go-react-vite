import * as React from "react";
import {
  Link,
  Outlet,
  createRootRoute,
  useLocation,
} from "@tanstack/react-router";

export const Route = createRootRoute({
  component: RootComponent,
  errorComponent: ({ error }) => `Error: ${error.message}`,
  notFoundComponent: NotFound,
});

function NotFound() {
  const loc = useLocation();
  return (
    <div
      style={{
        fontSize: "1.8rem",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      <h2>Not Found</h2>
      <pre>{loc.pathname}</pre>
    </div>
  );
}

function RootComponent() {
  return (
    <React.Fragment>
      <div
        style={{
          marginLeft: "auto",
          marginRight: "auto",
          maxWidth: "768px",
          border: "1px solid blue",
          padding: "2rem",
        }}
      >
        <div
          style={{
            display: "flex",
            // justifyContent: "center",
            alignItems: "center",
            gap: ".875rem",
          }}
        >
          <Link to="/">Home</Link>
          <Link to="/about">About</Link>
        </div>
        <hr style={{ marginBottom: "1.2rem" }} />
        <Outlet />
      </div>
    </React.Fragment>
  );
}
