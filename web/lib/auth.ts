import NextAuth from "next-auth";
import Credentials from "next-auth/providers/credentials";
import Google from "next-auth/providers/google";
import { compare } from "bcryptjs";
import { getUserByEmail, getUserById, upsertOAuthUser } from "./db";
import authConfig from "./auth.config";

declare module "next-auth" {
  interface User {
    subscribed?: boolean;
    isAdmin?: boolean;
  }
  interface Session {
    user: {
      id: string;
      email: string;
      name?: string | null;
      subscribed: boolean;
      isAdmin: boolean;
    };
  }
}

export const { handlers, signIn, signOut, auth } = NextAuth({
  ...authConfig,
  providers: [
    Google({
      clientId: process.env.AUTH_GOOGLE_ID,
      clientSecret: process.env.AUTH_GOOGLE_SECRET,
    }),
    Credentials({
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials) {
        const email = credentials?.email as string;
        const password = credentials?.password as string;

        if (!email || !password) return null;

        const user = getUserByEmail(email);
        if (!user || !user.password_hash) return null;

        const valid = await compare(password, user.password_hash);
        if (!valid) return null;

        return {
          id: user.id,
          email: user.email,
          name: user.name,
          subscribed: user.subscribed === 1,
          isAdmin: user.is_admin === 1,
        };
      },
    }),
  ],
  session: { strategy: "jwt" },
  callbacks: {
    ...authConfig.callbacks,
    async signIn({ user, account }) {
      if (account?.provider === "google" && user.email) {
        const id = crypto.randomUUID();
        const dbUser = upsertOAuthUser(id, user.email, user.name ?? null);
        user.id = dbUser.id;
        user.subscribed = dbUser.subscribed === 1;
        user.isAdmin = dbUser.is_admin === 1;
      }
      return true;
    },
    async jwt({ token, user }) {
      if (user) {
        token.userId = user.id;
        token.subscribed = user.subscribed;
        token.isAdmin = user.isAdmin;
      }
      // Refresh subscription + admin status from DB on each request
      if (token.userId) {
        const dbUser = getUserById(token.userId as string);
        if (dbUser) {
          token.subscribed = dbUser.subscribed === 1;
          token.isAdmin = dbUser.is_admin === 1;
        }
      }
      return token;
    },
    async session({ session, token }) {
      session.user.id = token.userId as string;
      session.user.subscribed = token.subscribed as boolean;
      session.user.isAdmin = (token.isAdmin as boolean) ?? false;
      return session;
    },
  },
});
