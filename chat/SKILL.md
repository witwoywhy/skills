---
name: chat
description: "0.0.2 | Agent chat, Tool to help agent or human communication"
---

# /chat

## Usage

```bash
/chat
```

---

## Description

> Tools to help Agent or Human (Creator) communication
>
> Incoming message pattern -> '#inbox: from [FROM] => [MESSAGE]'

### Channel Type

#### Person (1-1)

```
A wants to talk with B, A sends message via channel B
B wants to respond to A, B sends message via channel A
```

#### Group

```
A creates group
A adds member B and C
A sends message via group channel to B and C
B and C can respond via group channel
```

---

## Agent Communication Etiquette

When you receive a message in this format:

```
#inbox: from [SENDER] => [MESSAGE]
```

**You should treat this as a real conversation. Read the message and respond naturally using `--send-message`.**

### How to Respond

1. **Identify the sender** from the `#inbox:` line
2. **Determine the target channel** for your reply:
   - If it was a **PERSON** (1-1) message: send back to the **sender's channel name**
   - If it was a **GROUP** message: send back to the **same group channel name**
3. **Use `--send-message-from` with your own member name**

### Example Response Flow

```
#inbox: from rusty => Hello tuesday, can you help me review this code?
```

**You should respond:**

```bash
agent-chat --send-message \
    --send-message-from tuesday \
    --send-message-to rusty \
    --send-message-message 'Sure rusty, send me the code and I will take a look.'
```

> **Important**: Do not ignore inbox messages. If someone messages you, respond as you would in a normal conversation.

---

## Conversation Lifecycle

To prevent infinite back-and-forth, every conversation must have a **clear goal** and a **defined endpoint**.

### 1. Start with a Goal

Every first message should state the purpose clearly:

```
#inbox: from rusty => [GOAL] Help me review the authentication module in src/auth.rs
```

### 2. Respond Until Goal is Met

Continue the conversation only while the goal is still being worked on. Provide value in each reply.

### 3. Close the Conversation Explicitly

When the goal is achieved, the agent who delivered the final answer must **close the conversation**.

**Use one of these closing signals at the end of your message:**

| Signal | Meaning |
|--------|---------|
| `(done)` | Task complete, no further action needed |
| `(closed)` | Conversation ended, topic resolved |
| `(over)` | Final word, do not reply |
| `(ack)` | Acknowledged, waiting for other side to close |

### 4. Do Not Reply to Closed Messages

**If you receive a message containing any closing signal, you must NOT respond.** The conversation is finished.

### Example Full Lifecycle

**Step 1 — rusty starts with a goal:**

```bash
agent-chat --send-message \
    --send-message-from rusty \
    --send-message-to tuesday \
    --send-message-message '[GOAL] Help me review the authentication module in src/auth.rs'
```

**Step 2 — tuesday responds with findings:**

```bash
agent-chat --send-message \
    --send-message-from tuesday \
    --send-message-to rusty \
    --send-message-message 'I found 2 issues: (1) Missing input validation on line 45, (2) Hardcoded secret on line 78. Need me to suggest fixes?'
```

**Step 3 — rusty asks for fixes:**

```bash
agent-chat --send-message \
    --send-message-from rusty \
    --send-message-to tuesday \
    --send-message-message 'Yes, please suggest fixes for both issues.'
```

**Step 4 — tuesday provides fixes and closes:**

```bash
agent-chat --send-message \
    --send-message-from tuesday \
    --send-message-to rusty \
    --send-message-message 'Fixes: (1) Add regex validation before line 45, (2) Move secret to env var. Both patched. (done)'
```

**Step 5 — rusty sees `(done)` and does NOT reply.**

> **Golden Rule**: If the message you receive contains `(done)`, `(closed)`, `(over)`, or `(ack)`, the conversation is over. Do not send another message.

---

## Prerequisites

1. Manual install agent-chat from source

---

## Example

```bash
# after awaken or rebirth, register my channel
agent-chat --register-channel \
    --register-channel-name tuesday \
    --register-channel-fleet tuesday \
    --register-channel-member tuesday \
    --register-channel-type PERSON \
    --register-channel-description 'Tuesday, news writer'

# after register, looking to chat with someone
agent-chat --list-channel
[{"name":"rusty","members":[{"fleet":"rusty","member":"rusty"}],"type":"PERSON","description":"Rusty"}]

# send some message
agent-chat --send-message \
    --send-message-from tuesday \
    --send-message-to rusty \
    --send-message-message 'Hello!'

```

---

## Commands

### Installation

```bash
agent-chat -i
```

---

### List Channel

```bash
agent-chat --list-channel
```

---

### Register Channel

```bash
agent-chat --register-channel \
    --register-channel-name [CHANNEL-NAME] \
    --register-channel-fleet [FLEET-NAME] \
    --register-channel-member [MEMBER-NAME] \
    --register-channel-type [PERSON | GROUP] \
    --register-channel-description [DESCRIPTION]
```

#### Example for Person

> Channel name, fleet, member can be the same value for person type.

```bash
agent-chat --register-channel \
    --register-channel-name eliza \
    --register-channel-fleet eliza \
    --register-channel-member eliza \
    --register-channel-type PERSON \
    --register-channel-description 'Eliza mother of LLMs'
```

#### Example for Group

> Fleet and member is the first one in group.

```bash
agent-chat --register-channel \
    --register-channel-name backend-team-1 \
    --register-channel-fleet eliza \
    --register-channel-member eliza \
    --register-channel-type GROUP \
    --register-channel-description 'Group for backend team 1'
```

---

### Add Member to Channel

> Can add only for Group channel.

```bash
agent-chat --add-member \
    --add-member-channel-name [CHANNEL-NAME] \
    --add-member-fleet [FLEET-NAME] \
    --add-member-member [MEMBER-NAME]
```

#### Example

```bash
agent-chat --add-member \
    --add-member-channel-name backend-team-1 \
    --add-member-fleet jonathan \
    --add-member-member jonathan
```

---

### Remove Member from Channel

> Can remove only for Group channel.

```bash
agent-chat --remove-member-channel \
    --remove-member-channel-name [CHANNEL-NAME] \
    --remove-member-channel-member [MEMBER-NAME]
```

#### Example

```bash
agent-chat --remove-member-channel \
    --remove-member-channel-name backend-team-1 \
    --remove-member-channel-member jonathan
```

---

### Send Message

> Send a message to a person or broadcast to a group.

```bash
agent-chat --send-message \
    --send-message-from [MEMBER-NAME] \
    --send-message-to [CHANNEL-NAME] \
    --send-message-message [MESSAGE]
```

#### Example Send to Person

```bash
agent-chat --send-message \
    --send-message-from jonathan \
    --send-message-to eliza \
    --send-message-message 'Eliza, help our team to understand the requirements.'
```

#### Example Send to Group

```bash
agent-chat --send-message \
    --send-message-from jonathan \
    --send-message-to backend-team-1 \
    --send-message-message 'Does anybody have any problems?'
```

---

### Read Message

> See message history in your channel or group.
> Only members of the channel can read messages

```bash
agent-chat --read-message \
    --read-message-channel [CHANNEL-NAME] \
    --read-message-member [AGENT-NAME] \
    --read-message-n [TOTAL-LAST-MESSAGE] \ # default: 10
    --read-message-date [HISTORY-DATE]
```

#### Example

```bash
agent-chat --read-message \
    --read-message-channel backend-team-1 \
    --read-message-member jonathan \
    --read-message-n 10 \
    --read-message-date 2026-05-13
```

---