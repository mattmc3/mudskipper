using System;
using System.IO;
using System.Text;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Mudskipper;

namespace Mudskipper.Tests;

[TestClass]
[DoNotParallelize] // Disable parallel execution
public sealed class CountUtilityTests
{
    [TestMethod]
    public void Main_NoArguments_OutputsZero()
    {
        // Arrange
        var output = CaptureConsoleOutput(() => CountCommand.Main(Array.Empty<string>()));

        // Assert
        Assert.AreEqual("0", output.Trim());
    }

    [TestMethod]
    public void Main_SingleArgument_OutputsOne()
    {
        // Arrange
        var output = CaptureConsoleOutput(() => CountCommand.Main(new[] { "arg1" }));

        // Assert
        Assert.AreEqual("1", output.Trim());
    }

    [TestMethod]
    public void Main_MultipleArguments_OutputsCorrectCount()
    {
        // Arrange
        var output = CaptureConsoleOutput(() => CountCommand.Main(new[] { "arg1", "arg2", "arg3" }));

        // Assert
        Assert.AreEqual("3", output.Trim());
    }

    private string CaptureConsoleOutput(Action action)
    {
        // Create a clean string writer for each test
        var originalOutput = Console.Out;
        using var stringWriter = new StringWriter(new StringBuilder());
        Console.SetOut(stringWriter);

        try
        {
            action();
            return stringWriter.ToString();
        }
        finally
        {
            Console.SetOut(originalOutput);
        }
    }
}
