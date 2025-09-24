import 'dart:convert';
import 'dart:math';
import 'package:crypto/crypto.dart';
import 'package:shared_preferences/shared_preferences.dart';

class PasswordService {
  static const String _passwordHashKey = 'password_hash';
  static const String _saltKey = 'password_salt';
  static const String _isSetupKey = 'password_setup_complete';

  static PasswordService? _instance;
  static PasswordService get instance => _instance ??= PasswordService._();
  
  PasswordService._();

  late SharedPreferences _prefs;
  bool _initialized = false;

  Future<void> init() async {
    if (!_initialized) {
      _prefs = await SharedPreferences.getInstance();
      _initialized = true;
    }
  }

  // Check if password has been set up
  Future<bool> isPasswordSet() async {
    await init();
    return _prefs.getBool(_isSetupKey) ?? false;
  }

  // Generate a random salt
  String _generateSalt() {
    final random = Random.secure();
    final saltBytes = List<int>.generate(32, (i) => random.nextInt(256));
    return base64.encode(saltBytes);
  }

  // Hash password with salt
  String _hashPassword(String password, String salt) {
    final bytes = utf8.encode(password + salt);
    final digest = sha256.convert(bytes);
    return digest.toString();
  }

  // Set up password for the first time
  Future<bool> setupPassword(String password) async {
    await init();
    
    if (password.isEmpty || password.length < 4) {
      return false; // Password too short
    }

    try {
      final salt = _generateSalt();
      final hashedPassword = _hashPassword(password, salt);
      
      await _prefs.setString(_saltKey, salt);
      await _prefs.setString(_passwordHashKey, hashedPassword);
      await _prefs.setBool(_isSetupKey, true);
      
      return true;
    } catch (e) {
      return false;
    }
  }

  // Verify password
  Future<bool> verifyPassword(String password) async {
    await init();
    
    final salt = _prefs.getString(_saltKey);
    final storedHash = _prefs.getString(_passwordHashKey);
    
    if (salt == null || storedHash == null) {
      return false;
    }

    final hashedPassword = _hashPassword(password, salt);
    return hashedPassword == storedHash;
  }

  // Change password
  Future<bool> changePassword(String oldPassword, String newPassword) async {
    await init();
    
    // Verify old password first
    if (!await verifyPassword(oldPassword)) {
      return false;
    }

    if (newPassword.isEmpty || newPassword.length < 4) {
      return false; // New password too short
    }

    try {
      final salt = _generateSalt();
      final hashedPassword = _hashPassword(newPassword, salt);
      
      await _prefs.setString(_saltKey, salt);
      await _prefs.setString(_passwordHashKey, hashedPassword);
      
      return true;
    } catch (e) {
      return false;
    }
  }

  // Reset password system (for development/testing)
  Future<void> resetPasswordSystem() async {
    await init();
    await _prefs.remove(_passwordHashKey);
    await _prefs.remove(_saltKey);
    await _prefs.remove(_isSetupKey);
  }
}