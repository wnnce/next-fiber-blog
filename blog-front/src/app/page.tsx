import RichImage from '@/components/RichImage'

export default function Home() {
  const nums: number[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
  return (
    <div className="min-h-screen w-full">
      <div>
        <RichImage src="/b-oss/images/2024/0704/7d8e2b2e0aa20f139b8788aa38066403.webp" />
        <RichImage src="/b-oss/images/2024/0704/7d8e2b2e0aa20f139b8788aa3806640311.webp" />
      </div>
      {
        nums.map(num => {
          return (
            <div key={num} className="animate-on-scroll w-100 h-100 bg-gray">
              {num}
            </div>
          )
        })
      }
    </div>
  )
}