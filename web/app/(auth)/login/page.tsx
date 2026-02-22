import { Suspense } from "react";
import AuthForm from "../AuthForm";

export const metadata = {
  title: "Login - Lark",
};

export default function LoginPage() {
  return (
    <Suspense>
      <AuthForm mode="login" />
    </Suspense>
  );
}
