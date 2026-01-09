# Publishing PromptForge to VS Code Marketplace

## Prerequisites

1. **Azure DevOps Account** (free)
   - Sign up at: https://dev.azure.com

2. **VS Code Extension Manager (vsce)**
   ```bash
   npm install -g @vscode/vsce
   ```

## Step-by-Step Publishing

### Step 1: Create a Publisher

1. Go to: https://marketplace.visualstudio.com/manage
2. Sign in with your Microsoft/Azure account
3. Click "Create Publisher"
4. Fill in:
   - **Publisher ID**: `promptforge` (or your choice)
   - **Publisher Name**: `PromptForge`
   - **Description**: Brief description
5. Click "Create"

### Step 2: Get Personal Access Token

1. Go to: https://dev.azure.com
2. Click your profile → **Security**
3. Click **Personal Access Tokens**
4. Click **New Token**
5. Fill in:
   - **Name**: `VS Code Marketplace`
   - **Organization**: Select your org
   - **Expiration**: Choose duration
   - **Scopes**: Check **Marketplace (Manage)**
6. Click **Create**
7. **Copy the token** (you won't see it again!)

### Step 3: Login to vsce

```bash
cd vscode-extension
vsce login promptforge
```

When prompted, paste your Personal Access Token.

### Step 4: Publish

```bash
vsce publish
```

This will:
- Package the extension
- Upload to marketplace
- Make it available for download

## Alternative: Manual Upload

If you prefer manual upload:

```bash
# 1. Package the extension
cd vscode-extension
vsce package

# 2. Go to marketplace
# https://marketplace.visualstudio.com/manage

# 3. Click "New Extension" → "Visual Studio Code"

# 4. Upload the .vsix file
# (promptforge-editor-1.0.0.vsix)
```

## After Publishing

1. **Extension will be available** in VS Code Marketplace
2. **Users can install** via:
   - VS Code Extensions view → Search "PromptForge"
   - Command: `code --install-extension promptforge.promptforge-editor`
   - Marketplace URL: `https://marketplace.visualstudio.com/items?itemName=promptforge.promptforge-editor`

## Updating the Extension

To publish updates:

1. Update version in `package.json`
2. Run: `vsce publish`
3. New version appears in marketplace automatically

## Requirements Checklist

- ✅ `package.json` with proper metadata
- ✅ `LICENSE` file (MIT)
- ✅ `README.md` for marketplace
- ✅ Extension code works
- ✅ No errors in packaging

## Current Status

Your extension is ready to publish! Just follow the steps above.

