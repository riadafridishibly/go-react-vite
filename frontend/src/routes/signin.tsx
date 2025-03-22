import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/signin")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<div>
			<a href="/api/auth/google">Sign in with Google</a>
		</div>
	);
}
