import { StrictMode, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";

import { routeTree } from "./routeTree.gen";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { AuthContext, useAuth, type User } from "./auth";

const router = createRouter({
	routeTree,
	defaultErrorComponent: () => "Error Occured!",
	context: {
		auth: {
			user: null,
			isAuthenticated: false,
		},
	},
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

// Render the app
// biome-ignore lint/style/noNonNullAssertion: <explanation>
const rootElement = document.getElementById("root")!;
if (!rootElement.innerHTML) {
	const root = createRoot(rootElement);
	root.render(
		<StrictMode>
			<AuthProvider>
				<InnerApp />
			</AuthProvider>
		</StrictMode>,
	);
}

function InnerApp() {
	const auth = useAuth();
	const { user = null, isAuthenticated = false } = auth ?? {};

	const contextData = useMemo(
		() => ({
			auth: {
				user,
				isAuthenticated,
			},
		}),
		[user, isAuthenticated],
	);

	return <RouterProvider router={router} context={contextData} />;
}

function AuthProvider({ children }: { children: React.ReactNode }) {
	const [user, setUser] = useState<User | null>(null);
	const [isAuthenticated, setIsAuthenticated] = useState(false);

	return (
		<AuthContext.Provider
			value={{ user, setUser, isAuthenticated, setIsAuthenticated }}
		>
			{children}
		</AuthContext.Provider>
	);
}
