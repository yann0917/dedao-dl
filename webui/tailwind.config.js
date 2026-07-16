var config = {
    darkMode: ["class"],
    content: ["./index.html", "./src/**/*.{ts,tsx}"],
    theme: {
        extend: {
            colors: {
                border: "hsl(var(--border))",
                input: "hsl(var(--input))",
                ring: "hsl(var(--ring))",
                background: "hsl(var(--background))",
                foreground: "hsl(var(--foreground))",
                surface: {
                    page: "hsl(var(--bg-page))",
                    panel: "hsl(var(--bg-panel))",
                    soft: "hsl(var(--bg-soft))",
                    inverse: "hsl(var(--bg-inverse))",
                },
                text: {
                    primary: "hsl(var(--text-primary))",
                    secondary: "hsl(var(--text-secondary))",
                    muted: "hsl(var(--text-muted))",
                    inverse: "hsl(var(--text-inverse))",
                },
                primary: {
                    DEFAULT: "hsl(var(--primary))",
                    foreground: "hsl(var(--primary-foreground))",
                },
                accent: {
                    DEFAULT: "hsl(var(--accent))",
                    foreground: "hsl(var(--accent-foreground))",
                    soft: "hsl(var(--accent-soft))",
                },
                danger: {
                    DEFAULT: "hsl(var(--danger))",
                    foreground: "hsl(var(--danger-foreground))",
                    soft: "hsl(var(--danger-soft))",
                },
                success: {
                    DEFAULT: "hsl(var(--success))",
                    foreground: "hsl(var(--success-foreground))",
                    soft: "hsl(var(--success-soft))",
                },
                warning: {
                    DEFAULT: "hsl(var(--warning))",
                    foreground: "hsl(var(--warning-foreground))",
                    soft: "hsl(var(--warning-soft))",
                },
                secondary: {
                    DEFAULT: "hsl(var(--secondary))",
                    foreground: "hsl(var(--secondary-foreground))",
                },
                muted: {
                    DEFAULT: "hsl(var(--muted))",
                    foreground: "hsl(var(--muted-foreground))",
                },
                card: {
                    DEFAULT: "hsl(var(--card))",
                    foreground: "hsl(var(--card-foreground))",
                },
            },
            borderRadius: {
                lg: "var(--radius)",
                md: "calc(var(--radius) - 2px)",
                sm: "calc(var(--radius) - 4px)",
            },
            boxShadow: {
                soft: "0 24px 60px rgba(15, 23, 42, 0.12)",
            },
        },
    },
    plugins: [],
};
export default config;
