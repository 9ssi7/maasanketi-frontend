import { useEffect } from "react";

type Props = {
  onLogin?: (token: string) => void;
};

declare global {
  interface Window {
    google: {
      accounts: {
        id: {
          initialize: (options: unknown) => void;
          renderButton: (element: HTMLElement, options: unknown) => void;
          prompt: () => void;
        };
      };
    };
  }
}

export default function OAuthGoogleButton({ onLogin }: Props) {
  useEffect(() => {
    setTimeout(() => {
      function handleCredentialResponse(response: { credential: string }) {
        onLogin?.(response.credential);
      }
      window.google.accounts.id.initialize({
        client_id: import.meta.env.PUBLIC_GOOGLE_OAUTH_CLIENT_ID,
        callback: handleCredentialResponse,
      });

      window.google.accounts.id.renderButton(
        document.getElementById("g_id_sign")!,
        { theme: "outline" }
      );

      window.google.accounts.id.prompt(); // also display the One Tap dialog
    }, 1000);
  }, [onLogin]);
  return <div id="g_id_sign" className="w-full lg:w-auto"></div>;
}
