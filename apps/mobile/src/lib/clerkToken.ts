type ClerkTokenGetter = () => Promise<string | null>;

let clerkTokenGetter: ClerkTokenGetter | null = null;

export function setClerkTokenGetter(getter: ClerkTokenGetter | null): void {
  clerkTokenGetter = getter;
}

export async function getClerkToken(): Promise<string | null> {
  if (!clerkTokenGetter) {
    return null;
  }

  return clerkTokenGetter();
}
