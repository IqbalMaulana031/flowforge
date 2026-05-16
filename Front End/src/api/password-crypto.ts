const publicKeyB64 = import.meta.env.VITE_PASSWORD_ENCRYPTION_PUBLIC_KEY_B64;

function base64ToBytes(value: string) {
  const binary = atob(value);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i += 1) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}

function bytesToBase64(value: ArrayBuffer) {
  const bytes = new Uint8Array(value);
  let binary = '';
  for (const byte of bytes) {
    binary += String.fromCharCode(byte);
  }
  return btoa(binary);
}

async function importPasswordEncryptionKey() {
  if (!publicKeyB64) {
    throw new Error('Password encryption public key is not configured');
  }

  return crypto.subtle.importKey(
    'spki',
    base64ToBytes(publicKeyB64),
    {name: 'RSA-OAEP', hash: 'SHA-256'},
    false,
    ['encrypt'],
  );
}

export async function encryptPassword(password: string) {
  const key = await importPasswordEncryptionKey();
  const encrypted = await crypto.subtle.encrypt(
    {name: 'RSA-OAEP'},
    key,
    new TextEncoder().encode(password),
  );
  return bytesToBase64(encrypted);
}
