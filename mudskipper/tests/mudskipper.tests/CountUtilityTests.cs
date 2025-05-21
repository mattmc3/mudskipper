using NUnit.Framework;

namespace Mudskipper.Tests
{
    public class CountUtilityTests
    {
        [Test]
        public void Main_NoArguments_OutputsZero()
        {
            var result = CountUtility.Main(new string[] { });
            Assert.AreEqual(0, result);
        }

        [Test]
        public void Main_SingleArgument_OutputsOne()
        {
            var result = CountUtility.Main(new string[] { "arg1" });
            Assert.AreEqual(1, result);
        }

        [Test]
        public void Main_MultipleArguments_OutputsCorrectCount()
        {
            var result = CountUtility.Main(new string[] { "arg1", "arg2", "arg3" });
            Assert.AreEqual(3, result);
        }
    }
}