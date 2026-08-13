# MeetSpace

> A real-time collaborative development workspace where developers can code, communicate, execute code, and collaborate visually — all from one place.

MeetSpace is a full-stack collaborative coding platform built around a **Go backend**, **React frontend**, **PostgreSQL**, **Piston code execution**, and a **Pion-based WebRTC SFU**.

The goal is to combine the core capabilities developers normally use across multiple tools — collaborative coding, file management, real-time communication, code execution, and whiteboarding — into a single workspace.

---

## ✨ Features

### 🔐 Authentication

- User registration and login
- JWT-based authentication
- Secure password hashing
- Protected API routes
- Environment-based JWT secret configuration

### 👥 Collaborative Projects

- Create and manage projects
- Project membership
- Role-based access control (RBAC)
- Owner / Editor / Viewer roles
- Project invitations
- Member management

### 📁 Collaborative File Workspace

- Create files and folders
- Nested directory structure
- Rename files and folders
- Delete files and folders
- Persistent file contents
- Project-based file isolation

### 💻 Online Code Editor

- Monaco Editor integration
- JavaScript, Python and Go support
- File-based editing
- Code execution directly from the workspace
- Execution output and error handling

### ⚡ Code Execution

MeetSpace integrates with **Piston** to safely execute user-submitted code in isolated environments.

Supported runtimes include:

- JavaScript / Node.js
- Python
- Go

The Go backend communicates with Piston through an HTTP API rather than executing Docker commands directly.

### 🔄 Real-Time Collaboration

- WebSocket-based collaboration
- Real-time editor synchronization
- User presence
- Project-specific collaboration channels
- Real-time project updates

### 🎨 Collaborative Whiteboard

- Konva-based whiteboard
- Shape creation and manipulation
- Real-time synchronization
- Persistent whiteboard state
- Project-specific whiteboard data

### 🎙️ Voice & Video Collaboration

MeetSpace includes a custom **Pion WebRTC SFU** for real-time audio/video communication.

Features include:

- Voice communication
- Video communication
- WebRTC signaling
- Peer connection management
- Real-time media tracks
- Project-based communication rooms

---

# 🏗️ Architecture

```text
                         MeetSpace
                             │
              ┌──────────────┴──────────────┐
              │                             │
          React + Vite                  Go Backend
              │                             │
              │                    ┌────────┼─────────┐
              │                    │        │         │
              │                    ▼        ▼         ▼
              │               PostgreSQL  Piston    WebSocket
              │                                      Hub
              │                                        │
              │                                        ▼
              │                                  Pion SFU
              │                                        │
              │                                        ▼
              │                                   WebRTC
              │
              ├── Monaco Editor
              ├── File Explorer
              ├── Whiteboard
              └── Collaboration UI
