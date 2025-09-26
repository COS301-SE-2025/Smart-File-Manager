import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../constants.dart';

class SettingsService extends ChangeNotifier {
  static const String _colorPresetKey = 'graph_color_preset';
  static const String _namingConventionKey = 'folder_naming_convention';

  static SettingsService? _instance;
  static SettingsService get instance => _instance ??= SettingsService._();

  SettingsService._();

  late SharedPreferences _prefs;
  bool _initialized = false;

  Future<void> init() async {
    if (!_initialized) {
      _prefs = await SharedPreferences.getInstance();
      _initialized = true;
    }
  }

  // Color Preset Methods
  Future<void> setColorPreset(String preset) async {
    await init();
    await _prefs.setString(_colorPresetKey, preset);
    notifyListeners();
  }

  Future<String> getColorPreset() async {
    await init();
    return _prefs.getString(_colorPresetKey) ?? 'Default';
  }

  // Naming Convention Methods
  Future<void> setNamingConvention(String convention) async {
    await init();
    await _prefs.setString(_namingConventionKey, convention);
    notifyListeners();
  }

  Future<String> getNamingConvention() async {
    await init();
    return _prefs.getString(_namingConventionKey) ?? 'CAMEL';
  }

  // Color Preset Definitions
  static List<Color> getPresetColors(String preset) {
    switch (preset) {
      case "Default":
        return [
          kprimaryColor,
          const Color(0xff3B82F6),
          const Color(0xff10B981),
          const Color(0xffE74C3C),
          const Color(0xff9B59B6),
        ];
      case "Ocean":
        return [
          const Color(0xff0EA5E9),
          const Color(0xff06B6D4),
          const Color(0xff3B82F6),
          const Color(0xff1D4ED8),
          const Color(0xff0C4A6E),
        ];
      case "Forest":
        return [
          const Color(0xff10B981),
          const Color(0xff059669),
          const Color(0xff065F46),
          const Color(0xff047857),
          const Color(0xff064E3B),
        ];
      case "Sunset":
        return [
          const Color(0xffF59E0B),
          const Color(0xffEF4444),
          const Color(0xffDC2626),
          const Color(0xffB91C1C),
          const Color(0xff991B1B),
        ];
      case "Monochrome":
        return [
          const Color(0xff6B7280),
          const Color(0xff4B5563),
          const Color(0xff374151),
          const Color(0xff1F2937),
          const Color(0xff111827),
        ];
      default:
        return [
          kprimaryColor,
          const Color(0xff3B82F6),
          const Color(0xff10B981),
          const Color(0xffE74C3C),
          const Color(0xff9B59B6),
        ];
    }
  }

  static const List<String> availablePresets = [
    "Default",
    "Ocean",
    "Forest",
    "Sunset",
    "Monochrome",
  ];
  static const List<String> availableNamingConventions = [
    "CAMEL",
    "SNAKE",
    "PASCAL",
    "KEBAB",
    "SPACE",
  ];
}
