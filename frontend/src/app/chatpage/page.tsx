"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import ChatInterface from "@/components/ChatInterface";

function ChatPage() {
  const searchParams = useSearchParams();
  const threadId = searchParams.get("threadId");
  const acceleratorId = searchParams.get("acceleratorId");

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-l from-cyan-500 via-teal-600 to-green-700 dark:to-green-950 dark:via-teal-950 dark:from-gray-900">
      <div className="container mx-auto p-4">
        {threadId && acceleratorId ? (
          <ChatInterface threadId={threadId} acceleratorId={acceleratorId} />
        ) : (
          <p>Loading...</p>
        )}
      </div>
    </div>
  );
}

export default function MyWrapperComponent() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <ChatPage/>
    </Suspense>
  );
}
