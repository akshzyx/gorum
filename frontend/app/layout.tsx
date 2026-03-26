import "./globals.css";

export const metadata = {
  title: "Gorum",
  description: "Go run for it",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <link
          rel="stylesheet"
          href="https://font-awesome-pro-v7.vercel.app/css/fontawesome.css"
        />
        <link
          rel="stylesheet"
          href="https://font-awesome-pro-v7.vercel.app/css/allmain.css"
        />
      </head>
      <body>{children}</body>
    </html>
  );
}
