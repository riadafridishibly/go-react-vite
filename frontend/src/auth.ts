import { createContext, useContext } from "react";

export interface User {
	id: string;
	name: string;
	email: string;
	avatar: string;
}

interface AuthContextData {
	user: User | null;
	setUser: (user: User | null) => void;
	isAuthenticated: boolean;
	setIsAuthenticated: (isAuthenticated: boolean) => void;
}

export const AuthContext = createContext<AuthContextData | null>(null);

export const useAuth = () => {
	return useContext(AuthContext);
};
