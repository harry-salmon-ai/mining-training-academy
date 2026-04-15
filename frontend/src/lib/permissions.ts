type Role = "SUPERADMIN" | "ADMIN" | "INSTRUCTOR" | "SUPERVISOR" | "LEARNER";

const PERMISSIONS: Record<Role, Record<string, string[]>> = {
  SUPERADMIN: {
    modules: ["create", "read", "update", "delete", "publish", "archive"],
    users: ["create", "read", "update", "delete", "assign_role"],
    categories: ["create", "read", "update", "delete"],
    enrollments: ["create", "read", "update", "delete", "assign"],
    reports: ["read", "export"],
    settings: ["read", "update"],
    audit: ["read"],
  },
  ADMIN: {
    modules: ["create", "read", "update", "delete", "publish"],
    users: ["create", "read", "update"],
    categories: ["create", "read", "update"],
    enrollments: ["create", "read", "update", "assign"],
    reports: ["read", "export"],
    settings: ["read"],
    audit: ["read"],
  },
  INSTRUCTOR: {
    modules: ["create", "read", "update"],
    users: ["read"],
    categories: ["read"],
    enrollments: ["read", "assign"],
    reports: ["read"],
  },
  SUPERVISOR: {
    modules: ["read"],
    users: ["read"],
    enrollments: ["read", "assign"],
    reports: ["read"],
  },
  LEARNER: {
    modules: ["read"],
    enrollments: ["read"],
  },
};

export function canAccess(
  role: Role | undefined,
  resource: string,
  action: string
): boolean {
  if (!role) return false;
  return PERMISSIONS[role]?.[resource]?.includes(action) ?? false;
}
