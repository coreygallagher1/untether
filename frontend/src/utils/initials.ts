// Utility function to generate initials from user data
export const getInitials = (user: { first_name?: string; last_name?: string; username?: string }): string => {
  // If we have first and last name, use those
  if (user.first_name && user.last_name) {
    return (user.first_name.charAt(0) + user.last_name.charAt(0)).toUpperCase();
  }
  
  // If we only have first name, use first two characters
  if (user.first_name && !user.last_name) {
    if (user.first_name.length >= 2) {
      return user.first_name.substring(0, 2).toUpperCase();
    }
    return user.first_name.charAt(0).toUpperCase();
  }
  
  // If we only have last name, use first two characters
  if (!user.first_name && user.last_name) {
    if (user.last_name.length >= 2) {
      return user.last_name.substring(0, 2).toUpperCase();
    }
    return user.last_name.charAt(0).toUpperCase();
  }
  
  // Fallback to username logic if no first/last name
  const username = user.username || '';
  if (!username) return 'U';
  
  // Try to split on common separators first (like john.doe, john_doe, john-doe)
  const parts = username.split(/[._-]/);
  if (parts.length >= 2) {
    return (parts[0].charAt(0) + parts[1].charAt(0)).toUpperCase();
  }
  
  // Try to find camelCase patterns (like johnDoe)
  const camelCaseMatch = username.match(/^([a-z]+)([A-Z][a-z]*)/);
  if (camelCaseMatch) {
    return (camelCaseMatch[1].charAt(0) + camelCaseMatch[2].charAt(0)).toUpperCase();
  }
  
  // Try to find patterns with numbers (like john123, john123doe)
  const numberMatch = username.match(/^([a-zA-Z]+)(\d+)([a-zA-Z]*)/);
  if (numberMatch && numberMatch[3]) {
    return (numberMatch[1].charAt(0) + numberMatch[3].charAt(0)).toUpperCase();
  }
  
  // If it's a single word, try to find vowels to split on
  const vowels = /[aeiouAEIOU]/;
  const vowelIndex = username.search(vowels);
  if (vowelIndex > 0 && vowelIndex < username.length - 1) {
    return (username.charAt(0) + username.charAt(vowelIndex + 1)).toUpperCase();
  }
  
  // Fallback: take first character and last character if different
  if (username.length >= 2 && username.charAt(0) !== username.charAt(username.length - 1)) {
    return (username.charAt(0) + username.charAt(username.length - 1)).toUpperCase();
  }
  
  // Final fallback: take first two characters
  if (username.length >= 2) {
    return username.substring(0, 2).toUpperCase();
  }
  
  // Single character
  return username.charAt(0).toUpperCase();
};
