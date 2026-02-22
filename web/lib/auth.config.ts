import type { NextAuthConfig } from "next-auth";

// This config is safe for Edge Runtime (no Node.js-only imports).
// The full auth config in auth.ts extends this with Credentials provider.
export default {
  pages: {
    signIn: "/login",
  },
  callbacks: {
    authorized({ auth, request: { nextUrl } }) {
      const isProtected =
        nextUrl.pathname.startsWith("/play") ||
        nextUrl.pathname.startsWith("/admin");
      if (isProtected && !auth) return false;
      return true;
    },
  },
  providers: [], // Providers added in auth.ts
} satisfies NextAuthConfig;
