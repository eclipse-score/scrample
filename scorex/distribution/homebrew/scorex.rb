class Scorex < Formula
  desc "CLI for creating S-CORE skeleton projects"
  homepage "https://github.com/eclipse-score/score_scrample"
  version "0.1.0"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/eclipse-score/score_scrample/releases/download/v#{version}/scorex-#{version}-macos-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_ARM64"
    else
      url "https://github.com/eclipse-score/score_scrample/releases/download/v#{version}/scorex-#{version}-macos-x86_64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_X86_64"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/eclipse-score/score_scrample/releases/download/v#{version}/scorex-#{version}-linux-x86_64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX"
    end
  end

  def install
    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "scorex-macos-arm64" => "scorex"
      else
        bin.install "scorex-macos-x86_64" => "scorex"
      end
    elsif OS.linux?
      bin.install "scorex-linux-x86_64" => "scorex"
    end
  end

  test do
    system "#{bin}/scorex", "version"
  end
end
