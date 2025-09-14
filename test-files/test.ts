// Interface definitions
interface User {
  id: number;
  name: string;
  email: string;
  role: "admin" | "user" | "guest";
  settings?: UserSettings;
  metadata: Record<string, unknown>;
}

interface UserSettings {
  theme: "light" | "dark" | "system";
  notifications: boolean;
  language: string;
  timezone: string;
}

// Types
type FilterFunction<T> = (item: T, index: number, array: T[]) => boolean;
type SortDirection = "asc" | "desc";
type ValidationResult = { valid: boolean; errors: string[] };

// Utility functions
function validateEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

function formatDate(date: Date, format: string = "yyyy-mm-dd"): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");

  let result = format;
  result = result.replace("yyyy", year.toString());
  result = result.replace("mm", month);
  result = result.replace("dd", day);

  return result;
}

// Class implementation
class UserManager {
  private users: Map<number, User>;
  private nextId: number;
  private logger: (message: string) => void;

  constructor(initialUsers: User[] = [], logger?: (message: string) => void) {
    this.users = new Map();
    this.nextId = 1;
    this.logger = logger || console.log;

    // Initialize with any provided users
    if (initialUsers && initialUsers.length > 0) {
      for (const user of initialUsers) {
        if (user.id >= this.nextId) {
          this.nextId = user.id + 1;
        }
        this.users.set(user.id, this.sanitizeUser(user));
      }
    }
  }

  private sanitizeUser(user: User): User {
    const sanitized = { ...user };

    // Ensure email is valid
    if (sanitized.email && !validateEmail(sanitized.email)) {
      this.logger(`Invalid email for user ${user.id}: ${user.email}`);
      sanitized.email = "";
    }

    // Set default settings if not provided
    if (!sanitized.settings) {
      sanitized.settings = {
        theme: "system",
        notifications: true,
        language: "en",
        timezone: "UTC",
      };
    }

    return sanitized;
  }

  public addUser(user: Omit<User, "id">): User {
    const newUser: User = {
      ...user,
      id: this.nextId++,
    };

    const sanitizedUser = this.sanitizeUser(newUser);
    this.users.set(sanitizedUser.id, sanitizedUser);
    this.logger(`User added: ${sanitizedUser.id} - ${sanitizedUser.name}`);

    return sanitizedUser;
  }

  public updateUser(id: number, updates: Partial<Omit<User, "id">>): User | null {
    const existingUser = this.users.get(id);

    if (!existingUser) {
      this.logger(`Update failed: User ${id} not found`);
      return null;
    }

    const updatedUser: User = {
      ...existingUser,
      ...updates,
    };

    const sanitizedUser = this.sanitizeUser(updatedUser);
    this.users.set(id, sanitizedUser);
    this.logger(`User updated: ${sanitizedUser.id} - ${sanitizedUser.name}`);

    return sanitizedUser;
  }

  public deleteUser(id: number): boolean {
    if (!this.users.has(id)) {
      this.logger(`Delete failed: User ${id} not found`);
      return false;
    }

    this.users.delete(id);
    this.logger(`User deleted: ${id}`);
    return true;
  }

  public getUser(id: number): User | null {
    return this.users.get(id) || null;
  }

  public searchUsers(searchTerm: string): User[] {
    const results: User[] = [];
    const lowerSearchTerm = searchTerm.toLowerCase();

    for (const user of this.users.values()) {
      const nameMatch = user.name.toLowerCase().includes(lowerSearchTerm);
      const emailMatch = user.email.toLowerCase().includes(lowerSearchTerm);

      if (nameMatch || emailMatch) {
        results.push(user);
      }
    }

    return results;
  }

  public validateUser(user: User): ValidationResult {
    const errors: string[] = [];

    // Validate required fields
    if (!user.name || user.name.trim() === "") {
      errors.push("Name is required");
    } else if (user.name.length < 2) {
      errors.push("Name must be at least 2 characters");
    }

    if (!user.email || user.email.trim() === "") {
      errors.push("Email is required");
    } else if (!validateEmail(user.email)) {
      errors.push("Email format is invalid");
    }

    // Validate role
    if (!["admin", "user", "guest"].includes(user.role)) {
      errors.push("Role must be one of: admin, user, guest");
    }

    // Validate settings if present
    if (user.settings) {
      if (!["light", "dark", "system"].includes(user.settings.theme)) {
        errors.push("Theme must be one of: light, dark, system");
      }

      if (typeof user.settings.notifications !== "boolean") {
        errors.push("Notifications setting must be a boolean");
      }
    }

    return {
      valid: errors.length === 0,
      errors,
    };
  }

  public filterAndSortUsers<K extends keyof User>(
    filterFn: FilterFunction<User> | null = null,
    sortKey: K | null = null,
    sortDirection: SortDirection = "asc",
  ): User[] {
    // Convert map to array
    let result = Array.from(this.users.values());

    // Apply filter if provided
    if (filterFn) {
      result = result.filter(filterFn);
    }

    // Apply sorting if key provided
    if (sortKey) {
      result.sort((a, b) => {
        const valueA = a[sortKey];
        const valueB = b[sortKey];

        // Handle string comparisons
        if (typeof valueA === "string" && typeof valueB === "string") {
          return sortDirection === "asc" ? valueA.localeCompare(valueB) : valueB.localeCompare(valueA);
        }

        // Handle number comparisons
        if (typeof valueA === "number" && typeof valueB === "number") {
          return sortDirection === "asc" ? valueA - valueB : valueB - valueA;
        }

        // Default comparison for other types
        return 0;
      });
    }

    return result;
  }

  public exportUsers(format: "json" | "csv" = "json"): string {
    const users = Array.from(this.users.values());

    if (format === "json") {
      return JSON.stringify(users, null, 2);
    } else if (format === "csv") {
      // Create CSV header
      let csv = "id,name,email,role\n";

      // Add each user as a row
      for (const user of users) {
        csv += `${user.id},"${user.name}","${user.email}",${user.role}\n`;
      }

      return csv;
    } else {
      throw new Error(`Unsupported export format: ${format}`);
    }
  }

  public getUsersWithRoleCount(): Record<string, number> {
    const roleCounts: Record<string, number> = {
      admin: 0,
      user: 0,
      guest: 0,
    };

    for (const user of this.users.values()) {
      roleCounts[user.role]++;
    }

    return roleCounts;
  }

  public getTotalUsers(): number {
    return this.users.size;
  }
}

// Example usage
function demonstrateUserManager(): void {
  // Create a custom logger
  const logger = (message: string): void => {
    const timestamp = new Date().toISOString();
    console.log(`[${timestamp}] ${message}`);
  };

  // Initialize with some users
  const initialUsers: User[] = [
    {
      id: 101,
      name: "Alice Johnson",
      email: "alice@example.com",
      role: "admin",
      metadata: { lastLogin: "2023-05-10", loginCount: 42 },
    },
    {
      id: 102,
      name: "Bob Smith",
      email: "bob@example.com",
      role: "user",
      settings: {
        theme: "dark",
        notifications: false,
        language: "en",
        timezone: "America/New_York",
      },
      metadata: { lastLogin: "2023-05-09", loginCount: 17 },
    },
  ];

  const userManager = new UserManager(initialUsers, logger);

  // Add a new user
  const newUser = userManager.addUser({
    name: "Charlie Davis",
    email: "charlie@example.com",
    role: "guest",
    metadata: { lastLogin: null, loginCount: 0 },
  });

  // Update a user
  if (newUser) {
    const updated = userManager.updateUser(newUser.id, {
      role: "user",
      settings: {
        theme: "light",
        notifications: true,
        language: "fr",
        timezone: "Europe/Paris",
      },
    });

    if (updated) {
      logger(`User role changed from 'guest' to '${updated.role}'`);
    }
  }

  // Search for users
  const searchResults = userManager.searchUsers("alice");
  logger(`Search found ${searchResults.length} users`);

  // Filter and sort users
  const activeAdmins = userManager.filterAndSortUsers((user) => user.role === "admin", "name", "asc");

  logger(`Found ${activeAdmins.length} active admins`);

  // Export users
  const jsonExport = userManager.exportUsers("json");
  logger(`Exported ${userManager.getTotalUsers()} users to JSON`);

  // Get role statistics
  const roleCounts = userManager.getUsersWithRoleCount();
  for (const [role, count] of Object.entries(roleCounts)) {
    logger(`Role ${role}: ${count} users`);
  }

  // Complex filtering with multiple conditions
  const complexFilter = userManager.filterAndSortUsers(
    (user) => {
      // Check if user has dark theme
      const hasDarkTheme = user.settings?.theme === "dark";

      // Check if user has logged in recently (metadata check)
      const lastLogin = user.metadata?.lastLogin as string;
      const recentLogin = lastLogin ? new Date(lastLogin) > new Date("2023-05-01") : false;

      // Check if user has significant login count
      const loginCount = user.metadata?.loginCount as number;
      const highLoginCount = typeof loginCount === "number" && loginCount > 10;

      // Combined conditions
      return hasDarkTheme || (recentLogin && highLoginCount);
    },
    "id",
    "desc",
  );

  logger(`Complex filter found ${complexFilter.length} matching users`);
}

// Run the demonstration
demonstrateUserManager();

// Export for potential reuse
export { User, UserSettings, UserManager, validateEmail, formatDate };
